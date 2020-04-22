package rclone

import (
	"fmt"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/maputils"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/stringutils"
	"github.com/sirupsen/logrus"
	"go/types"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

/* Struct */

type RemoteServiceAccounts struct {
	RemoteEnvVar    string
	ServiceAccounts []pathutils.Path
}

type RemoteServiceAccount struct {
	RemoteEnvVar       string
	ServiceAccountPath string
}

type ServiceAccountManager struct {
	log                         *logrus.Entry
	remoteServiceAccountFolders map[string]string
	remoteServiceAccounts       map[string]RemoteServiceAccounts
}

/* Public */

func NewServiceAccountManager(serviceAccountFolders map[string]string) *ServiceAccountManager {
	return &ServiceAccountManager{
		log:                         logger.GetLogger("sa_manager"),
		remoteServiceAccountFolders: serviceAccountFolders,
		remoteServiceAccounts:       make(map[string]RemoteServiceAccounts),
	}
}

func (m *ServiceAccountManager) LoadServiceAccounts(remotePaths []string) error {
	m.log.Trace("Loading service accounts")

	// iterate remotes
	for _, remotePath := range remotePaths {
		// ignore junk paths
		if remotePath == "" {
			continue
		}

		// parse remote name and retrieve folder
		remoteName := stringutils.FromLeftUntil(remotePath, ":")
		remoteServiceAccountFolder, err := maputils.GetStringMapValue(m.remoteServiceAccountFolders, remoteName,
			false)
		if err != nil {
			m.log.Tracef("Service account folder was not found for: %q, skipping...", remoteName)
			continue
		}

		// service accounts loaded for this remote?
		if _, ok := m.remoteServiceAccounts[remoteName]; ok {
			continue
		}

		// retrieve service files
		serviceAccountFiles, _ := pathutils.GetPathsInFolder(remoteServiceAccountFolder, true,
			false, func(path string) *string {
				lowerPath := strings.ToLower(path)

				// ignore non json files
				if !strings.HasSuffix(lowerPath, ".json") {
					return nil
				}

				return &path
			})

		// were service accounts found?
		if len(serviceAccountFiles) == 0 {
			m.log.Tracef("No service accounts found for %q in: %v", remoteName, remoteServiceAccountFolder)
			continue
		}

		// sort service files
		re := regexp.MustCompile("[0-9]+")
		sort.SliceStable(serviceAccountFiles, func(i, j int) bool {
			is := stringutils.NewOrExisting(re.FindString(serviceAccountFiles[i].RealPath), "0")
			js := stringutils.NewOrExisting(re.FindString(serviceAccountFiles[j].RealPath), "0")

			in, err := strconv.Atoi(is)
			if err != nil {
				return false
			}
			jn, err := strconv.Atoi(js)
			if err != nil {
				return false
			}

			return in < jn
		})

		// add to remote service accounts var
		v := RemoteServiceAccounts{
			RemoteEnvVar:    fmt.Sprintf("RCLONE_CONFIG_%s_SERVICE_ACCOUNT_FILE", strings.ToUpper(remoteName)),
			ServiceAccounts: serviceAccountFiles,
		}
		m.remoteServiceAccounts[remoteName] = v

		m.log.Debugf("Loaded %d service accounts for remote %q (env: %v)", len(serviceAccountFiles), remoteName,
			v.RemoteEnvVar)
	}

	return nil
}

func (m *ServiceAccountManager) GetServiceAccount(remotePaths ...string) ([]*RemoteServiceAccount, error) {
	var serviceAccounts []*RemoteServiceAccount
	var err error
	successfulRemotes := make(map[string]*types.Nil)

	for _, remotePath := range remotePaths {
		saFound := false

		// parse remote name
		remoteName := stringutils.FromLeftUntil(remotePath, ":")
		if remoteName == "" {
			// no remote name was parsed, so ignore this request
			m.log.Tracef("No remote determined for: %q, not providing service account", remotePath)
			continue
		}

		// service accounts loaded for this remote?
		remote, ok := m.remoteServiceAccounts[remoteName]
		if !ok || len(remote.ServiceAccounts) == 0 {
			// no service accounts found for this remote
			m.log.Tracef("No service accounts loaded for remote: %q, not providing service account", remoteName)
			continue
		}

		// have we already set a service account for this remote?
		if _, ok := successfulRemotes[strings.ToLower(remoteName)]; ok {
			continue
		}

		// find unbanned service account
		for _, sa := range remote.ServiceAccounts {
			// does the cache already contain this service account?
			if exists, _ := cache.IsBanned(sa.RealPath); exists {
				// service account is currently banned
				continue
			}

			// this service account is unbanned
			serviceAccounts = append(serviceAccounts, &RemoteServiceAccount{
				RemoteEnvVar:       remote.RemoteEnvVar,
				ServiceAccountPath: sa.RealPath,
			})

			saFound = true
			break
		}

		if saFound {
			// we found a service account, check for next remote
			successfulRemotes[strings.ToLower(remoteName)] = nil
			continue
		}

		// if we are here, no more service accounts were available
		m.log.Warnf("No more service accounts available for remote: %q", remoteName)
		err = fmt.Errorf("failed finding available service account for remote: %q", remoteName)
		break
	}

	return serviceAccounts, err
}

func (m *ServiceAccountManager) ServiceAccountsCount() int {
	n := 0

	for _, remote := range m.remoteServiceAccounts {
		n += len(remote.ServiceAccounts)
	}

	return n
}

func AnyRemotesBanned(remotes []string) (bool, time.Time) {
	var banned bool
	var expires time.Time

	// ignore empty remotes slice
	if remotes == nil {
		return banned, expires
	}

	// format remotes into remote names if possible
	var checkRemotes []string
	for _, remote := range remotes {
		checkRemotes = append(checkRemotes, stringutils.FromLeftUntil(remote, ":"))
	}

	// iterate remotes
	for _, remote := range checkRemotes {
		banned, expires = cache.IsBanned(remote)
		if banned {
			break
		}
	}

	return banned, expires
}

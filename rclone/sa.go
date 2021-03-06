package rclone

import (
	"fmt"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/maputils"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/reutils"
	"github.com/l3uddz/crop/stringutils"
	"github.com/sirupsen/logrus"
	"go/types"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	remoteServiceAccountFolders map[string][]string
	remoteServiceAccounts       map[string]RemoteServiceAccounts
	parallelism                 int
}

var (
	mtx    sync.Mutex
	mcache *ttlcache.Cache
)

/* Private */

func init() {
	mcache = ttlcache.NewCache()
	_ = mcache.SetTTL(60 * time.Minute)
	mcache.SetExpirationCallback(mcacheItemExpired)
}

func mcacheItemExpired(key string, _ interface{}) {
	log.Debugf("Cleared SA from mcache: %s", key)
}

func addServiceAccountsToTempCache(serviceAccounts []*RemoteServiceAccount) {
	for _, sa := range serviceAccounts {
		_ = mcache.Set(sa.ServiceAccountPath, nil)
	}
}

/* Public */

func NewServiceAccountManager(serviceAccountFolders map[string][]string, parallelism int) *ServiceAccountManager {
	return &ServiceAccountManager{
		log:                         logger.GetLogger("sa_manager"),
		remoteServiceAccountFolders: serviceAccountFolders,
		remoteServiceAccounts:       make(map[string]RemoteServiceAccounts),
		parallelism:                 parallelism,
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

		// parse remote name and retrieve service account folder(s)
		remoteName := stringutils.FromLeftUntil(remotePath, ":")
		remoteServiceAccountFolders, err := maputils.GetStringKeysBySliceValue(m.remoteServiceAccountFolders, remoteName)
		if err != nil {
			m.log.Tracef("Service account folder(s) not found for: %q, skipping...", remoteName)
			continue
		}

		// service accounts loaded for this remote?
		if _, ok := m.remoteServiceAccounts[remoteName]; ok {
			continue
		}

		// load service account files in all folders
		totalServiceAccountFiles := make([]pathutils.Path, 0)

		for _, remoteServiceAccountFolder := range remoteServiceAccountFolders {
			// retrieve service files within this folder
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
			sort.SliceStable(serviceAccountFiles, func(i, j int) bool {
				is := reutils.GetEveryNumber(serviceAccountFiles[i].RealPath)
				js := reutils.GetEveryNumber(serviceAccountFiles[j].RealPath)

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

			totalServiceAccountFiles = append(totalServiceAccountFiles, serviceAccountFiles...)
		}

		// were service accounts found?
		if len(totalServiceAccountFiles) == 0 {
			m.log.Tracef("No service accounts found for %q in: %v", remoteName, remoteServiceAccountFolders)
			continue
		}

		// add to remote service accounts var
		v := RemoteServiceAccounts{
			RemoteEnvVar:    ConfigToEnv(remoteName, "SERVICE_ACCOUNT_FILE"),
			ServiceAccounts: totalServiceAccountFiles,
		}
		m.remoteServiceAccounts[remoteName] = v

		m.log.Debugf("Loaded %d service accounts for remote %q (env: %v)", len(totalServiceAccountFiles),
			remoteName, v.RemoteEnvVar)
	}

	return nil
}

func (m *ServiceAccountManager) GetRandomServiceAccount(remotePath string) (string, error) {
	// parse remote name
	remoteName := stringutils.FromLeftUntil(remotePath, ":")
	if remoteName == "" {
		// no remote name was parsed, so ignore this request
		m.log.Tracef("No remote determined for: %q, not providing service account", remotePath)
		return "", nil
	}

	// service accounts loaded for this remote?
	remote, ok := m.remoteServiceAccounts[remoteName]
	if !ok || len(remote.ServiceAccounts) == 0 {
		// no service accounts found for this remote
		m.log.Tracef("No service accounts loaded for remote: %q, not providing service account", remoteName)
		return "", nil
	}

	// random service account
	rand.Seed(time.Now().Unix())
	sa := remote.ServiceAccounts[rand.Intn(len(remote.ServiceAccounts))]

	return sa.RealPath, nil
}

func (m *ServiceAccountManager) GetServiceAccount(remotePaths ...string) ([]*RemoteServiceAccount, error) {
	var serviceAccounts []*RemoteServiceAccount
	var err error
	successfulRemotes := make(map[string]*types.Nil)

	// acquire global lock
	mtx.Lock()
	defer mtx.Unlock()

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

			// has this service account been issued within N seconds?
			if _, err := mcache.Get(sa.RealPath); err == nil {
				// this sa was in our memory cache and has not expired yet
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

	// were service accounts found?
	if err == nil && m.parallelism > 1 && len(serviceAccounts) > 0 {
		// there may be multiple routines requesting service accounts
		// attempt to prevent service account from being re-used (unless explicitly removed by a successful operation)
		addServiceAccountsToTempCache(serviceAccounts)
	}

	return serviceAccounts, err
}

func (m *ServiceAccountManager) ServiceAccountsCount() int {
	n := 0
	t := make(map[string]int)

	for _, remote := range m.remoteServiceAccounts {
		for _, sa := range remote.ServiceAccounts {
			if _, ok := t[sa.RealPath]; !ok {
				t[sa.RealPath] = 0
				n++
			}
		}
	}

	return n
}

func RemoveServiceAccountsFromTempCache(serviceAccounts []*RemoteServiceAccount) {
	mtx.Lock()
	defer mtx.Unlock()

	for _, sa := range serviceAccounts {
		_ = mcache.Remove(sa.ServiceAccountPath)
	}
}

func AnyRemotesBanned(remotes []string) (bool, time.Time) {
	var banned bool
	var expires time.Time

	// ignore empty remotes slice
	if remotes == nil {
		return banned, expires
	}

	// format remotes into remote names if possible
	checkRemotes := make([]string, 0)
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

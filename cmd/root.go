package cmd

import (
	"fmt"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/logger"
	"github.com/l3uddz/crop/pathutils"
	"github.com/l3uddz/crop/rclone"
	"github.com/l3uddz/crop/runtime"
	"github.com/l3uddz/crop/stringutils"
	"github.com/nightlyone/lockfile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

var (
	// Global flags
	flagLogLevel     = 0
	flagConfigFolder = pathutils.GetDefaultConfigPath()
	flagConfigFile   = "config.yaml"
	flagCachePath    = "cache"
	flagLogFile      = "activity.log"
	flagLockFile     = "crop.lock"
	flagDryRun       bool
	flagNoDedupe     bool

	// Global command specific
	flagUploader string

	// Global vars
	log   *logrus.Entry
	flock lockfile.Lockfile
)

var rootCmd = &cobra.Command{
	Use:   "crop",
	Short: "CLI application to assist harvesting your media",
	Long: `A CLI application that can be used to harvest your local media.
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Parse persistent flags
	rootCmd.PersistentFlags().StringVar(&flagConfigFolder, "config-dir", flagConfigFolder, "Config folder")
	rootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", flagConfigFile, "Config file")
	rootCmd.PersistentFlags().StringVarP(&flagCachePath, "cache", "d", flagCachePath, "Cache path")
	rootCmd.PersistentFlags().StringVarP(&flagLogFile, "log", "l", flagLogFile, "Log file")
	rootCmd.PersistentFlags().StringVarP(&flagLockFile, "lock", "f", flagLockFile, "Lock file")
	rootCmd.PersistentFlags().CountVarP(&flagLogLevel, "verbose", "v", "Verbose level")

	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "Dry run mode")
}

func initCore(showAppInfo bool) {
	// Set core variables
	if !rootCmd.PersistentFlags().Changed("config") {
		flagConfigFile = filepath.Join(flagConfigFolder, flagConfigFile)
	}
	if !rootCmd.PersistentFlags().Changed("cache") {
		flagCachePath = filepath.Join(flagConfigFolder, flagCachePath)
	}
	if !rootCmd.PersistentFlags().Changed("log") {
		flagLogFile = filepath.Join(flagConfigFolder, flagLogFile)
	}
	if !rootCmd.PersistentFlags().Changed("lock") {
		flagLockFile = filepath.Join(flagConfigFolder, flagLockFile)
	}

	// Init Logging
	if err := logger.Init(flagLogLevel, flagLogFile); err != nil {
		log.WithError(err).Fatal("Failed to initialize logging")
	}

	log = logger.GetLogger("crop")

	// Init File Lock
	if err := acquireFileLock(); err != nil {
		log.WithError(err).Fatalf("Failed acquiring file lock for %q", flagLockFile)
	}

	// Init Config
	if err := config.Init(flagConfigFile); err != nil {
		log.WithError(err).Fatal("Failed to initialize config")
	}

	setConfigOverrides()

	// Init Cache
	if err := cache.Init(flagCachePath, flagLogLevel); err != nil {
		log.WithError(err).Fatal("Failed to initialize cache")
	}

	// Init Rclone
	if err := rclone.Init(config.Config); err != nil {
		log.WithError(err).Fatal("Failed to initialize rclone")
	}

	// Show App Info
	if showAppInfo {
		showUsing()
	}
}

func setConfigOverrides() {
	// set dry-run if enabled by flag
	if flagDryRun {
		config.Config.Rclone.DryRun = true
	}
}

func acquireFileLock() error {
	f, err := lockfile.New(flagLockFile)
	if err != nil {
		return err
	}

	flock = f

	// loop until lock has been acquired
	for {
		err = flock.TryLock()
		switch {
		case err == nil:
			// lock has been acquired
			return nil
		case err == lockfile.ErrBusy:
			// another instance is already running
			log.Warnf("There is another crop instance running, re-checking in 1 minute...")
			time.Sleep(1 * time.Minute)
		default:
			// an un-expected error, propagate down-stream
			return err
		}
	}
}

func releaseFileLock() {
	if err := flock.Unlock(); err != nil {
		log.WithError(err).Fatalf("Failed releasing file lock for %q", flagLockFile)
	}
}

func showUsing() {
	// show app info
	log.Infof("Using %s = %s (%s@%s)", stringutils.LeftJust("VERSION", " ", 10),
		runtime.Version, runtime.GitCommit, runtime.Timestamp)
	logger.ShowUsing()
	config.ShowUsing()
	cache.ShowUsing()
	log.Info("------------------")
}

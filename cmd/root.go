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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	// Global flags
	flagLogLevel     = 0
	flagConfigFolder = pathutils.GetCurrentBinaryPath()
	flagConfigFile   = "config.yaml"
	flagCacheFile    = "cache.json"
	flagLogFile      = "activity.log"

	flagDryRun bool

	// Global vars
	log *logrus.Entry
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
	rootCmd.PersistentFlags().StringVarP(&flagCacheFile, "cache", "d", flagCacheFile, "Cache file")
	rootCmd.PersistentFlags().StringVarP(&flagLogFile, "log", "l", flagLogFile, "Log file")
	rootCmd.PersistentFlags().CountVarP(&flagLogLevel, "verbose", "v", "Verbose level")

	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "Dry run mode")
}

func initCore(showAppInfo bool) {
	// Set core variables
	if !rootCmd.PersistentFlags().Changed("config") {
		flagConfigFile = filepath.Join(flagConfigFolder, flagConfigFile)
	}
	if !rootCmd.PersistentFlags().Changed("cache") {
		flagCacheFile = filepath.Join(flagConfigFolder, flagCacheFile)
	}
	if !rootCmd.PersistentFlags().Changed("log") {
		flagLogFile = filepath.Join(flagConfigFolder, flagLogFile)
	}

	// Init Logging
	if err := logger.Init(flagLogLevel, flagLogFile); err != nil {
		log.WithError(err).Fatal("Failed to initialize logging")
	}

	log = logger.GetLogger("crop")

	// Init Config
	if err := config.Init(flagConfigFile); err != nil {
		log.WithError(err).Fatal("Failed to initialize config")
	}

	// Init Cache
	if err := cache.Init(flagCacheFile); err != nil {
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

func showUsing() {
	// show app info
	log.Infof("Using %s = %s (%s@%s)", stringutils.LeftJust("VERSION", " ", 10),
		runtime.Version, runtime.GitCommit, runtime.Timestamp)
	logger.ShowUsing()
	config.ShowUsing()
	cache.ShowUsing()
	log.Info("------------------")
}

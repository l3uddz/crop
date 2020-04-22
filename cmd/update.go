package cmd

import (
	"bufio"
	"github.com/blang/semver"
	"github.com/l3uddz/crop/cache"
	"github.com/l3uddz/crop/runtime"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"os"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to latest version",
	Long:  `This command can be used to self-update to the latest version.`,

	Run: func(cmd *cobra.Command, args []string) {
		// init core
		initCore(false)
		defer cache.Close()

		// parse current version
		v, err := semver.Parse(runtime.Version)
		if err != nil {
			log.WithError(err).Fatal("Failed parsing current build version")
		}

		// detect latest version
		log.Info("Checking for the latest version...")
		latest, found, err := selfupdate.DetectLatest("l3uddz/crop")
		if err != nil {
			log.WithError(err).Fatal("Failed determining latest available version")
		}

		// check version
		if !found || latest.Version.LTE(v) {
			log.Infof("Already using the latest version: %v", runtime.Version)
			return
		}

		// ask update
		log.Infof("Do you want to update to the latest version: %v? (y/n):", latest.Version)
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil || (input != "y\n" && input != "n\n") {
			log.Fatal("Failed validating input...")
		} else if input == "n\n" {
			return
		}

		// get existing executable path
		exe, err := os.Executable()
		if err != nil {
			log.WithError(err).Fatal("Failed locating current executable path")
		}

		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			log.WithError(err).Fatal("Failed updating existing binary to latest release")
		}

		log.Infof("Successfully updated to the latest version: %v", latest.Version)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

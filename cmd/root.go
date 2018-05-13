package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Region string

var rootCmd = &cobra.Command{
	Long: `Uses the Aws go-sdk to provision and launch an ec2 instance.
Checks for a default VPC and will create it if its missing.
The instance is configured using user-data and will install docker.
Once docker is installed, the webapp will be pulled and ran as a container.`,
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&Region, "region", "r", "eu-west-1", "AWS region to use.")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

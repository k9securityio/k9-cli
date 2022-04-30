/*
Copyright © 2022 The K9CLI Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cmd contains all cobra commands
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/k9securityio/k9-cli/core"
)

// TODO take in path to sync with, select date to pull for, perform the sync

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync your local database with a report delivered to your AWS account.",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, _ := cmd.Flags().GetString(`bucket`)
		customerID, _ := cmd.Flags().GetString(`customer_id`)
		accountID, _ := cmd.Flags().GetString(`account`)
		// reportHome, _ := cmd.Flags().GetString(`report-home`)
		concurrency, _ := cmd.Flags().GetInt(`concurrency`)
		verbose, _ := cmd.Flags().GetBool(`verbose`)
		dryrun, _ := cmd.Flags().GetBool(`dryrun`)
		stdout := cmd.OutOrStdout()
		stderr := cmd.ErrOrStderr()

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Fprintf(stderr, "Error retrieving AWS configuration: %v+\n", err)
			os.Exit(1)
			return
		}

		s3db, err := core.LoadS3DB(s3.NewFromConfig(cfg), bucket)
		if err != nil {
			fmt.Fprintf(stderr, "Error loading remote database: %v+\n", err)
			os.Exit(1)
			return
		}

		err = core.Sync(stdout, stderr, s3db,
			manager.NewDownloader(s3.NewFromConfig(cfg)),
			bucket, customerID, accountID, concurrency, dryrun, verbose)
		if err != nil {
			fmt.Fprintf(stderr, "%v+\n", err)
			os.Exit(1)
		}
	},
}

// init defines and wires flags
func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().String(`path`, `.`, `Local path where reports will be stored`)
	viper.BindPFlag(`local-path`, syncCmd.Flags().Lookup(`path`))

	syncCmd.Flags().String(`bucket`, ``, `S3 bucket location of your K9 secure inbox`)
	viper.BindPFlag(`bucket`, syncCmd.Flags().Lookup(`bucket`))
	syncCmd.MarkFlagRequired(`bucket`)

	syncCmd.Flags().String(`customer_id`, ``, `K9 customer ID reports to download`)
	viper.BindPFlag(`customer_id`, syncCmd.Flags().Lookup(`customer_id`))
	syncCmd.MarkFlagRequired(`customer_id`)

	syncCmd.Flags().Int(`concurrency`, 4, `number of concurrent downloads`)

	syncCmd.Flags().String(`account`, ``, `AWS account for which reports will be downloaded`)
	syncCmd.Flags().Bool(`dryrun`, false, `don't perform the download`)

	viper.BindPFlag(`account`, syncCmd.Flags().Lookup(`account`))

}

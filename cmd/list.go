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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/k9securityio/k9-cli/core"
)

// TODO take in path to list with, select date to pull for, perform the list

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List customers, accounts, or reports in a local or remote repository.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error retrieving AWS configuration: %v+\n", err)
			os.Exit(1)
		}
		err = core.List(
			os.Stdout,
			cfg,
			cmd.Flags().Lookup(`bucket`).Value.String(),
			cmd.Flags().Lookup(`customer_id`).Value.String(),
			cmd.Flags().Lookup(`account`).Value.String())
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error retrieving the qualified list of reports: %v+\n", err)
			os.Exit(1)
		}
	},
}

// init defines and wires flags
func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP(`local`, `l`, ``, `list local resources`)

	listCmd.Flags().String(`path`, `.`, `local path where reports will be stored`)
	viper.BindPFlag(`local-path`, listCmd.Flags().Lookup(`path`))

	listCmd.Flags().String(`bucket`, ``, `S3 bucket localtion of your K9 secure-inbox (required)`)
	listCmd.MarkFlagRequired(`bucket`)
	viper.BindPFlag(`bucket`, listCmd.Flags().Lookup(`bucket`))

	listCmd.Flags().String(`account`, ``, `AWS account for which reports will be downloaded`)
	viper.BindPFlag(`account`, listCmd.Flags().Lookup(`account`))

	listCmd.Flags().String(`customer_id`, ``, `K9 customer ID reports to download`)
	viper.BindPFlag(`customer_id`, listCmd.Flags().Lookup(`customer_id`))

}

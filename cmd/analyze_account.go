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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/k9securityio/k9-cli/core"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var analyzeAccountCmd = &cobra.Command{
	Use:   "account",
	Short: `Analyze the specified account`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error retrieving AWS configuration: %v+\n", err)
			os.Exit(1)
		}

		stdout := cmd.OutOrStdout()

		customerID, _ := cmd.Flags().GetString(FLAG_CUSTOMER_ID)
		accountID, _ := cmd.Flags().GetString(FLAG_ACCOUNT)
		apiHost, _ := cmd.Flags().GetString("api")
		if apiHost == "" {
			apiHost = "api.k9security.io"
		}

		fmt.Fprintf(stdout, "Starting analysis of %v account %v using %v\n", customerID, accountID, apiHost)
		err = core.AnalyzeAccount(os.Stdout, cfg, apiHost, customerID, accountID)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error triggering analysis for %v account %v: %v+\n", customerID, accountID, err)
			os.Exit(1)
		}
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeAccountCmd)
	analyzeAccountCmd.Flags().String(`account`, ``, "The AWS account number for analysis (required)")
	analyzeAccountCmd.MarkFlagRequired(`account`)
	viper.BindPFlag(`account`, analyzeAccountCmd.Flags().Lookup(`account`))

	analyzeAccountCmd.Flags().String(`customer_id`, ``, `K9 customer ID that owns the account`)
	analyzeAccountCmd.MarkFlagRequired(`customer_id`)
	viper.BindPFlag(`customer_id`, analyzeAccountCmd.Flags().Lookup(`customer_id`))

	analyzeAccountCmd.Flags().String(`api`, ``, `K9 API to use for analysis`)
	viper.BindPFlag(`api`, analyzeAccountCmd.Flags().Lookup(`api`))
}

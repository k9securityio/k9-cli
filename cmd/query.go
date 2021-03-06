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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Lookup the effect of last scanned access control configuration by kind",
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.PersistentFlags().String(FLAG_ANALYSIS_DATE, ``, `Use snapshot from the specified date in YYYY-MM-DD (required)`)

	queryCmd.PersistentFlags().String(FLAG_FORMAT, `json`, `Output format [csv|json] (default: json)`)
	viper.BindPFlag(`query_format`, queryResourceCmd.Flags().Lookup(FLAG_FORMAT))

	queryCmd.PersistentFlags().String(FLAG_CUSTOMER_ID, ``, `K9 customer ID for analysis (required)`)
	queryCmd.MarkPersistentFlagRequired(FLAG_CUSTOMER_ID)
	queryCmd.PersistentFlags().String(FLAG_ACCOUNT, ``, `AWS account ID for analysis (required)`)
	queryCmd.MarkPersistentFlagRequired(FLAG_ACCOUNT)
}

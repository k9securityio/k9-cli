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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// queryRisksCmd represents the risks command
var queryRisksCmd = &cobra.Command{
	Use:   "risks",
	Short: "Lookup various risks discovered in an analysis",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("risks called, output format: %v\n", cmd.Flags().Lookup(`format`).Value)
	},
}

func init() {
	queryCmd.AddCommand(queryRisksCmd)

	queryRisksCmd.PersistentFlags().String(`format`, `json`, `Output format as one of: [ json | csv | tap | pdf ]`)
	viper.BindPFlag(`query_format`, queryRisksCmd.Flags().Lookup(`format`))
	queryRisksCmd.PersistentFlags().String(`analysis-date`, ``,
		`Use snapshot from the specified date in YYYY-MM-DD (required)`)
	queryRisksCmd.MarkFlagRequired(`analysis-date`)

	queryRisksCmd.PersistentFlags().String(`customer_id`, ``, `K9 customer ID for analysis (required)`)
	queryRisksCmd.MarkFlagRequired(`customer_id`)
	queryRisksCmd.PersistentFlags().String(`account`, ``, `AWS account ID for analysis (required)`)
	queryRisksCmd.MarkFlagRequired(`account`)
}

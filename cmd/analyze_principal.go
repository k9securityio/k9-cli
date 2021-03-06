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

// analyzePrincipalCmd represents the principal command
var analyzePrincipalCmd = &cobra.Command{
	Use:   "principal",
	Short: "Analyze access for the specified principal",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("principal called")
		// client.AnalyzePrincipal(cmd.Flag(`account`).Value, cmd.Flag(`service`).Value)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzePrincipalCmd)

	analyzePrincipalCmd.Flags().String(`account`, ``, "The AWS account number for analysis (required)")
	analyzePrincipalCmd.MarkFlagRequired(`account`)
	viper.BindPFlag(`account`, analyzePrincipalCmd.Flags().Lookup(`account`))

	analyzePrincipalCmd.Flags().StringArray(`service`, []string{}, "A list of service names to evaluate")
}

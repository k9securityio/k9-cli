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

// queryResourceAccessCmd represents the resource-access command
var queryResourceAccessCmd = &cobra.Command{
	Use:   "resource-access",
	Short: "Lookup one or more resource access summaries",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("resource-access called")
	},
}

func init() {
	queryCmd.AddCommand(queryResourceAccessCmd)

	queryResourceAccessCmd.Flags().String("format", `json`, `Output format [csv|json] (default: json)`)
	viper.BindPFlag(`query_format`, queryResourceAccessCmd.Flags().Lookup(`format`))
	queryResourceAccessCmd.Flags().StringArray("resource", []string{}, `A list of resource ARNs to include`)
}

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
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// queryPrincipalAccessCmd represents the principal-access command
var queryPrincipalAccessCmd = &cobra.Command{
	Use:   "principal-access",
	Short: "Lookup one or more principal-access summary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("principal-access called")
	},
}

func init() {
	queryCmd.AddCommand(queryPrincipalAccessCmd)

	queryPrincipalAccessCmd.Flags().String("format", `json`, `Output format [csv|json] (default: json)`)
	queryPrincipalAccessCmd.Flags().StringArray("principal", []string{}, `A list of principals to include`)
}

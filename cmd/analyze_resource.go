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
)

// analyzeResourceCmd represents the resource command
var analyzeResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("analyze resource called")
	},
}

var ()

func init() {
	analyzeCmd.AddCommand(analyzeResourceCmd)
	analyzeResourceCmd.Flags().StringArray(`service`, []string{}, "A list of service names")
	analyzeResourceCmd.Flags().String(`resource-arn`, ``, "The resource to analyze (required)")
	analyzeResourceCmd.MarkFlagRequired(`resource-arn`)
	analyzeResourceCmd.Flags().String(`principal-arn`, ``, "The principal to analyze")
}

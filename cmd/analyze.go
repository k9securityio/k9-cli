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

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   `analyze`,
	Short: `Start analysis on current configuration`,
	Long:  `< need longer description with examples >`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify a target to analyze")
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

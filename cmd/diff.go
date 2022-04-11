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

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: `Calculate the difference between a snapshot and last scan`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("diff called")
	},
}

// init defines and wires flags
func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().String(`format`, `csv`, `Output format: [csv]`)
	viper.BindPFlag(`diff_format`, diffCmd.Flags().Lookup(`format`))
}

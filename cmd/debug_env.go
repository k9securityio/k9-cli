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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/k9securityio/k9-cli/core"
	"github.com/spf13/cobra"
)

// debugEnvCmd represents the debugEnv command
var debugEnvCmd = &cobra.Command{
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("debug-env called")
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Printf("Error retrieving AWS configuration: %+v\n", err)
		} else {
			//			fmt.Printf("AWS Configuration: %+v\n", cfg)
			//			fmt.Printf("AWS Credentials: %+v\n", cfg.Credentials)
		}

		bucket, _ := cmd.Flags().GetString(`bucket`)
		if len(bucket) > 0 {
			client := s3.NewFromConfig(cfg)
			s3db, err := core.LoadS3DB(client, bucket)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Unable to load s3 database, %v\n", err)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), `Bucket database summary:`)
				s3db.Dump(cmd.OutOrStdout(), true)
			}
		} else {
			db, err := core.LoadLocalDB(cmd.Flags().Lookup(`report-home`).Value.String())
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Unable to load local database, %v\n", err)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), `Local database summary:`)
				db.Dump(cmd.OutOrStdout(), true)
			}
		}
	},
	Use:   "debug-env",
	Short: "Display debugEnv information about the environment",
}

func init() {
	rootCmd.AddCommand(debugEnvCmd)
	debugEnvCmd.Flags().String(`bucket`, ``, `S3 bucket location of your K9 secure inbox`)

}

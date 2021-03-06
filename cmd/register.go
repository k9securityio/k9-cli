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

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register with k9",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("register called for %v at %v\n",
			cmd.Flag(`customer-name`).Value,
			cmd.Flag(`technical-contact-email`).Value)
		// client.Register(cmd.Flag(`customer-name).Value, cmd.Flag(`technical-contact-email`).Value)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().String("customer-name", ``,
		"A full name for the registering customer (required)")
	registerCmd.MarkFlagRequired(`customer-name`)
	viper.BindPFlag(`customer_name`, registerCmd.Flags().Lookup(`customer-name`))

	registerCmd.Flags().String("technical-contact-email", ``,
		"A valid email address for the customer's technical contact (required)")
	registerCmd.MarkFlagRequired(`technical-contact-email`)
	viper.BindPFlag(`technical_contact_email`, registerCmd.Flags().Lookup(`technical-contact-email`))
}

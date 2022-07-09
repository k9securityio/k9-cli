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
	"io"
	"os"
	"time"

	"github.com/k9securityio/k9-cli/core"
	"github.com/k9securityio/k9-cli/views"
	"github.com/spf13/cobra"
)

// queryRisksPrivilegeEscalationCmd represents the risks command
var queryRisksPrivilegeEscalationCmd = &cobra.Command{
	Use:     "privilege-escalation",
	Aliases: []string{`iam-admins`},
	Short:   "Show privilege escalation risks",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool(FLAG_VERBOSE)
		format, _ := cmd.Flags().GetString(FLAG_FORMAT)
		customerID, _ := cmd.Flags().GetString(FLAG_CUSTOMER_ID)
		accountID, _ := cmd.Flags().GetString(FLAG_ACCOUNT)
		analysisDate, _ := cmd.Flags().GetString(FLAG_ANALYSIS_DATE)
		reportHome, _ := cmd.Flags().GetString(FLAG_REPORT_HOME)
		stdout := cmd.OutOrStdout()
		stderr := cmd.ErrOrStderr()

		var reportDateTime *time.Time
		if len(analysisDate) > 0 {
			td, err := time.Parse(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT, analysisDate)
			if err != nil {
				fmt.Fprintf(stderr, "invalid analysis-date: %v\n", analysisDate)
			}
			reportDateTime = &td
		}

		DoQueryRisksPrivilegeEscalation(stdout, stderr, reportHome, customerID, accountID, format, reportDateTime, verbose)
	},
}

func init() {
	queryRisksCmd.AddCommand(queryRisksPrivilegeEscalationCmd)
}

// DoQueryRisksPrivilegeEscalation
func DoQueryRisksPrivilegeEscalation(stdout, stderr io.Writer, reportHome, customerID, accountID, format string, analysisDate *time.Time, verbose bool) {
	// load the local report database
	db, err := core.LoadLocalDB(reportHome)
	if err != nil {
		fmt.Printf("Unable to load local database, %v\n", err)
		os.Exit(1)
		return
	}
	if verbose {
		defer DumpDBStats(stderr, &db)
	}

	// determine the file name for the desired report
	path := db.GetPathForCustomerAccountTimeKind(customerID, accountID, analysisDate, core.REPORT_TYPE_PREFIX_PRINCIPALS)
	if path == nil || len(*path) <= 0 {
		fmt.Fprintf(stderr, "No report found for customer: %v account: %v date: %v\n", customerID, accountID, analysisDate)
		os.Exit(1)
		return
	}

	f, err := os.Open(*path)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the specified report: %v\n", err)
		os.Exit(1)
		return
	}

	// Open and load the report
	records, err := core.LoadPrincipalsReport(f)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to analyze the specified report: %v\n", err)
		os.Exit(1)
		return
	}

	// reducer - apply filtering or detective logic
	output := []core.PrincipalsReportItem{}
	for _, r := range records {
		if r.PrincipalIsIAMAdmin {
			output = append(output, r)
		}
	}
	views.Display(stdout, stderr, format, output)
}

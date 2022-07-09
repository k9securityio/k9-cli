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
	"encoding/json"
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
		verbose, err := cmd.Flags().GetBool(FLAG_VERBOSE)
		format, err := cmd.Flags().GetString(FLAG_FORMAT)
		customerID, err := cmd.Flags().GetString(FLAG_CUSTOMER_ID)
		accountID, err := cmd.Flags().GetString(FLAG_ACCOUNT)
		analysisDate, err := cmd.Flags().GetString(FLAG_ANALYSIS_DATE)
		reportHome, err := cmd.Flags().GetString(FLAG_REPORT_HOME)
		stdout := cmd.OutOrStdout()
		stderr := cmd.ErrOrStderr()

		if err != nil {
			fmt.Fprintf(stderr, err.Error())
			os.Exit(1)
			return
		}

		DoQueryRisksPrivilegeEscalation(stdout, stderr, reportHome, customerID, accountID, analysisDate, format, verbose)
	},
}

func init() {
	queryRisksCmd.AddCommand(queryRisksPrivilegeEscalationCmd)
}

func DoQueryRisksPrivilegeEscalation(stdout, stderr io.Writer, reportHome, customerID, accountID, analysisDate, format string, verbose bool) {
	// load the local report database
	db, err := core.LoadLocalDB(reportHome)
	if err != nil {
		fmt.Printf("Unable to load local database, %v\n", err)
		os.Exit(1)
		return
	}
	if verbose {
		defer func() {
			total, accounts, customers := db.Sizes()
			fmt.Fprintf(stderr,
				"Local database:\n\tCustomers:\t\t%v\n\tAccounts:\t\t%v\n\tTotal analysis dates: \t%v\n",
				customers, accounts, total)
		}()
	}

	// determine the file name for the desired report
	var path string
	if len(analysisDate) > 0 {
		reportDateTime, err := time.Parse(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT, analysisDate)
		if err != nil {
			fmt.Fprintf(stderr, "Invalid analysis-date: %v\n", analysisDate)
			os.Exit(1)
			return
		}
		if qr := db.GetPathForCustomerAccountTimeKind(
			customerID, accountID, reportDateTime,
			core.REPORT_TYPE_PREFIX_PRINCIPALS); qr != nil {
			path = *qr
		} else {
			fmt.Fprintf(stderr,
				"No such report: %v, %v, %v, total records: %v\n",
				customerID, accountID,
				reportDateTime.Format(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT),
				db.Size())
			os.Exit(1)
			return
		}
	} else {
		// latest
		fmt.Fprintln(stderr, `running latest report`)
		if qr := db.GetPathForCustomerAccountLatestKind(
			customerID, accountID, core.REPORT_TYPE_PREFIX_PRINCIPALS); qr != nil {
			path = *qr
		} else {
			fmt.Fprintf(stderr,
				"No such report: %v, %v, total records: %v\n",
				customerID, accountID, db.Size())
			os.Exit(1)
			return
		}
	}

	f, err := os.Open(path)
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

	// transform for output
	switch format {
	case `pdf`:
	case `csv`:
		views.WriteCSVTo(stdout, stderr, output)
	case `tap`:
	case `json`:
		b, err := json.Marshal(output)
		if err != nil {
			fmt.Fprintln(stderr, `unable to marshal report to json`)
		}
		fmt.Fprintln(stdout, string(b))
	default:
		fmt.Fprintln(stderr, `invalid output type`)
	}
}

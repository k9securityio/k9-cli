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

		verbose, _ := cmd.Flags().GetBool(`verbose`)
		customerID := cmd.Flags().Lookup(`customer_id`).Value.String()
		accountID := cmd.Flags().Lookup(`account`).Value.String()
		analysisDate, _ := cmd.Flags().GetString(`analysis-date`)

		// load the local report database
		db, err := core.LoadLocalDB(cmd.Flags().Lookup(`report-home`).Value.String())
		if err != nil {
			fmt.Printf("Unable to load local database, %v\n", err)
		}
		if verbose {
			defer func() {
				total, accounts, customers := db.Sizes()
				fmt.Fprintf(cmd.ErrOrStderr(),
					"Local database:\n\tCustomers:\t\t%v\n\tAccounts:\t\t%v\n\tTotal analysis dates: \t%v\n",
					customers, accounts, total)
			}()
		}

		// determine the file name for the desired report
		var path string
		if len(analysisDate) > 0 {
			reportDateTime, err := time.Parse(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT, analysisDate)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Invalid analysis-date: %v\n", analysisDate)
				os.Exit(1)
				return
			}
			if qr := db.GetPathForCustomerAccountTimeKind(
				customerID, accountID, reportDateTime,
				core.REPORT_TYPE_PREFIX_PRINCIPALS); qr != nil {
				path = *qr
			} else {
				fmt.Fprintf(cmd.ErrOrStderr(),
					"No such report: %v, %v, %v, total records: %v\n",
					customerID, accountID,
					reportDateTime.Format(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT),
					db.Size())
				os.Exit(1)
				return
			}
		} else {
			// latest
			fmt.Fprintln(cmd.ErrOrStderr(), `running latest report`)

			if qr := db.GetPathForCustomerAccountLatestKind(
				customerID, accountID, core.REPORT_TYPE_PREFIX_PRINCIPALS); qr != nil {
				path = *qr
			} else {
				// TODO handle no such report
				fmt.Fprintf(cmd.ErrOrStderr(),
					"No such report: %v, %v, total records: %v\n",
					customerID, accountID, db.Size())
				os.Exit(1)
				return
			}
		}

		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Unable to open the specified report: %v\n", err)
			os.Exit(1)
			return
		}

		// Open and load the report
		records, err := core.LoadPrincipalsReport(f)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Unable to analyze the specified report: %v\n", err)
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
		// branch on the PersistentFlag `format` - should be one of json, csv, tap, or PDF
		switch cmd.Flags().Lookup(`format`).Value.String() {
		case `pdf`:
		case `csv`:
			views.WriteCSVTo(os.Stdout, cmd.ErrOrStderr(), output)
		case `tap`:
		case `json`:
			b, err := json.Marshal(output)
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), `unable to marshal report to json`)
			}
			fmt.Println(string(b))
		default:
			fmt.Fprintln(cmd.ErrOrStderr(), `invalid output type`)
		}
	},
}

func init() {
	queryRisksCmd.AddCommand(queryRisksPrivilegeEscalationCmd)
}

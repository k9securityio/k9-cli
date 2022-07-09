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

// diffPrincipalsCmd represents the principals subcommand of diff
var diffPrincipalsCmd = &cobra.Command{
	Use:   "principals",
	Short: `Calculate the difference between a principals snapshot and last scan`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool(`verbose`)
		customerID, _ := cmd.Flags().GetString(`customer_id`)
		accountID, _ := cmd.Flags().GetString(`account`)
		analysisDate, _ := cmd.Flags().GetString(`analysis-date`)
		reportHome, _ := cmd.Flags().GetString(`report-home`)
		stdout := cmd.OutOrStdout()
		stderr := cmd.ErrOrStderr()

		if len(analysisDate) <= 0 {
			fmt.Fprintln(stderr, `an analysis-date is required for comparison`)
			os.Exit(1)
		}

		reportDateTime, err := time.Parse(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT, analysisDate)
		if err != nil {
			fmt.Fprintf(stderr, "invalid analysis-date: %v\n", analysisDate)
			os.Exit(1)
		}

		DoDiffPrincipals(stdout, stderr, reportHome, customerID, accountID, reportDateTime, verbose)
	},
}

// init defines and wires flags
func init() {
	diffCmd.AddCommand(diffPrincipalsCmd)
}

// DoDiffPrincipals
func DoDiffPrincipals(stdout, stderr io.Writer, reportHome, customerID, accountID string, analysisDate time.Time, verbose bool) {
	// load the local report database
	db, err := core.LoadLocalDB(reportHome)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to load local database, %v\n", err)
		os.Exit(1)
	}

	// get the latest analysis
	var latestReportPath, targetReportPath string

	if qr := db.GetPathForCustomerAccountTimeKind(
		customerID, accountID, nil, core.REPORT_TYPE_PREFIX_PRINCIPALS); qr != nil {
		latestReportPath = *qr
	} else {
		fmt.Fprintf(stderr,
			"No such latest report: %v, %v, total records: %v\n",
			customerID, accountID, db.Size())
		os.Exit(1)
	}

	// get the target analysis
	// determine the file name for the desired report
	if qr := db.GetPathForCustomerAccountTimeKind(
		customerID, accountID, &analysisDate,
		core.REPORT_TYPE_PREFIX_PRINCIPALS); qr != nil {
		targetReportPath = *qr
	} else {
		fmt.Fprintf(stderr,
			"No such target report: %v, %v, %v, total records: %v\n",
			customerID, accountID,
			analysisDate.Format(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT),
			db.Size())
		os.Exit(1)
	}

	// open and load the reports
	lf, err := os.Open(latestReportPath)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the latest report: %v\n", err)
		os.Exit(1)
	}
	tf, err := os.Open(targetReportPath)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the target report: %v\n", err)
		os.Exit(1)
	}

	latest := &core.PrincipalsReport{}
	err = core.LoadReport(lf, latest)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the latest report: %v\n", err)
		os.Exit(1)
	}
	target := &core.PrincipalsReport{}
	err = core.LoadReport(tf, target)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the target report: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Fprintf(stderr,
			"Target Analysis: %v, records: %v\nLatest Analysis: %v, records: %v\n",
			analysisDate, len(latest.Items), latest.Items[0].AnalysisTime, len(target.Items))
	}

	// index on principal ARN for each ReportItem
	targetByARN := map[string]core.PrincipalsReportItem{}
	for _, ri := range target.Items {
		targetByARN[ri.PrincipalARN] = ri
	}

	// This loop marks the ARNs that it sees in the latest report
	// and subsequently verifies that the target report does not
	// contain records that have yet to be seen.
	seen := map[string]struct{}{}
	mark := struct{}{}
	diffs := []core.PrincipalsReportItemDifference{}
	for _, ri := range latest.Items {
		seen[ri.PrincipalARN] = mark
		if ti, ok := targetByARN[ri.PrincipalARN]; !ok {
			diffs = append(diffs, ri.AddedDiff())
		} else if !ri.Equivalent(ti) {
			diffs = append(diffs, ri.Diff(ti))
		}
	}
	for _, ri := range target.Items {
		if _, ok := seen[ri.PrincipalARN]; !ok {
			diffs = append(diffs, ri.DeletedDiff())
		}
	}
	views.WriteCSVTo(stdout, stderr, diffs)
}

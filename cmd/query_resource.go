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

// queryResourceCmd represents the resource command
var queryResourceCmd = &cobra.Command{
	Use:     "resource",
	Aliases: []string{"resources"},
	Short:   "Lookup one or more resources",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool(FLAG_VERBOSE)
		format, _ := cmd.Flags().GetString(FLAG_FORMAT)
		customerID, _ := cmd.Flags().GetString(FLAG_CUSTOMER_ID)
		accountID, _ := cmd.Flags().GetString(FLAG_ACCOUNT)
		analysisDate, _ := cmd.Flags().GetString(FLAG_ANALYSIS_DATE)
		reportHome, _ := cmd.Flags().GetString(FLAG_REPORT_HOME)
		stdout := cmd.OutOrStdout()
		stderr := cmd.ErrOrStderr()
		arns, _ := cmd.Flags().GetStringSlice(FLAG_ARNS)
		names, _ := cmd.Flags().GetStringSlice(FLAG_NAMES)
		resourcesFilter := map[string]bool{}

		var reportDateTime *time.Time
		if len(analysisDate) > 0 {
			td, err := time.Parse(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT, analysisDate)
			if err != nil {
				fmt.Fprintf(stderr, "invalid analysis-date: %v\n", analysisDate)
				os.Exit(1)
				return
			}
			reportDateTime = &td
		}

		for _, p := range arns {
			resourcesFilter[p] = true
		}
		for _, p := range names {
			resourcesFilter[p] = true
		}

		DoQueryResource(stdout, stderr,
			reportHome, customerID, accountID, format,
			reportDateTime,
			verbose,
			resourcesFilter)

	},
}

func init() {
	queryCmd.AddCommand(queryResourceCmd)

	queryResourceCmd.Flags().StringSlice(FLAG_ARNS, []string{}, `A list of resource ARNs to retrieve`)
	queryResourceCmd.Flags().StringSlice(FLAG_NAMES, []string{}, `A list of resource names to retrieve`)
}

// DoQueryResource is the high-level query and filtering logic for querying resource reports. Externalized for testability.
func DoQueryResource(stdout, stderr io.Writer,
	reportHome, customerID, accountID, format string,
	analysisDate *time.Time,
	verbose bool,
	resources map[string]bool) {

	// load the local report database
	db, err := core.LoadLocalDB(reportHome)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to load local database, %v\n", err)
		os.Exit(1)
		return
	}

	if verbose {
		defer DumpDBStats(stderr, &db)
	}

	// determine the file name for the desired report
	path := db.GetPathForCustomerAccountTimeKind(customerID, accountID, analysisDate, core.REPORT_TYPE_PREFIX_RESOURCES)
	if path == nil || len(*path) <= 0 {
		fmt.Fprintf(stderr, "No report found for customer: %v account: %v date: %v\n", customerID, accountID, analysisDate)
		os.Exit(1)
		return
	}

	// get the report
	lf, err := os.Open(*path)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the requested report: %v\n", err)
		os.Exit(1)
		return
	}
	report := &core.ResourcesReport{}
	err = core.LoadReport(lf, report)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the requested report: %v\n", err)
		os.Exit(1)
		return
	}

	if verbose {
		fmt.Fprintf(stderr, "Target Analysis: %v, records: %v\n", analysisDate, len(report.Items))
	}

	if len(resources) <= 0 {
		views.Display(stdout, stderr, format, report.Items)
		return
	}

	results := []core.ResourcesReportItem{}
	for _, ri := range report.Items {
		if _, ok := resources[ri.ResourceARN]; ok {
			results = append(results, ri)
			continue
		}
		if _, ok := resources[ri.ResourceName]; ok {
			results = append(results, ri)
			continue
		}
	}
	views.Display(stdout, stderr, format, results)
}

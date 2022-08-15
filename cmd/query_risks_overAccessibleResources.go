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

// queryRisksOverAccessibleResourcesCmd represents the risks command
var queryRisksOverAccessibleResourcesCmd = &cobra.Command{
	Use:     "over-accessible-resources",
	Aliases: []string{`accessible`},
	Short:   "Show over accessible resource risks",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool(FLAG_VERBOSE)
		format, _ := cmd.Flags().GetString(FLAG_FORMAT)
		customerID, _ := cmd.Flags().GetString(FLAG_CUSTOMER_ID)
		accountID, _ := cmd.Flags().GetString(FLAG_ACCOUNT)
		analysisDate, _ := cmd.Flags().GetString(FLAG_ANALYSIS_DATE)
		reportHome, _ := cmd.Flags().GetString(FLAG_REPORT_HOME)
		stdout := cmd.OutOrStdout()
		stderr := cmd.ErrOrStderr()
		services, _ := cmd.Flags().GetStringSlice(FLAG_SERVICE)

		maxAdmins, _ := cmd.Flags().GetInt(FLAG_MAX_ADMIN)
		maxRead, _ := cmd.Flags().GetInt(FLAG_MAX_READ)
		maxWrite, _ := cmd.Flags().GetInt(FLAG_MAX_WRITE)
		maxDelete, _ := cmd.Flags().GetInt(FLAG_MAX_DELETE)

		policy := AccessibilityPolicy{
			AdminCap:  maxAdmins,
			ReadCap:   maxRead,
			WriteCap:  maxWrite,
			DeleteCap: maxDelete,
		}

		var reportDateTime *time.Time
		if len(analysisDate) > 0 {
			td, err := time.Parse(core.FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT, analysisDate)
			if err != nil {
				fmt.Fprintf(stderr, "invalid analysis-date: %v\n", analysisDate)
				os.Exit(1)
			}
			reportDateTime = &td
		}

		serviceMap := map[string]bool{}
		for _, s := range services {
			serviceMap[s] = true
		}

		DoQueryOverAccessibleResources(stdout, stderr,
			reportHome, customerID, accountID, format,
			reportDateTime,
			verbose,
			serviceMap,
			policy)
	},
}

func init() {
	queryRisksCmd.AddCommand(queryRisksOverAccessibleResourcesCmd)

	queryRisksOverAccessibleResourcesCmd.Flags().StringSlice(FLAG_SERVICE, []string{}, "A list of service names to evaluate")
	queryRisksOverAccessibleResourcesCmd.MarkFlagRequired(FLAG_SERVICE)

	queryRisksOverAccessibleResourcesCmd.Flags().Int(FLAG_MAX_ADMIN, 5, "The maximum number of principals with ADMIN access to a resource.")
	queryRisksOverAccessibleResourcesCmd.Flags().Int(FLAG_MAX_READ, 5, "The maximum number of principals with READ to a resource.")
	queryRisksOverAccessibleResourcesCmd.Flags().Int(FLAG_MAX_WRITE, 5, "The maximum number of principals with WRITE to a resource.")
	queryRisksOverAccessibleResourcesCmd.Flags().Int(FLAG_MAX_DELETE, 5, "The maximum number of principals with DELETE to a resource.")
}

func DoQueryOverAccessibleResources(stdout, stderr io.Writer,
	reportHome, customerID, accountID, format string,
	analysisDate *time.Time,
	verbose bool,
	services map[string]bool,
	policy AccessibilityPolicy) {

	// load the local report database
	db, err := core.LoadLocalDB(reportHome)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to load local database, %v\n", err)
		os.Exit(1)
	}

	if verbose {
		defer DumpDBStats(stderr, &db)
	}

	// determine the file name fo rthe desired report
	path := db.GetPathForCustomerAccountTimeKind(customerID, accountID, analysisDate, core.REPORT_TYPE_PREFIX_RESOURCE_ACCESS_SUMMARIES)
	if path == nil || len(*path) <= 0 {
		fmt.Fprintf(stderr, "No report found for customer: %v account: %v date: %v\n", customerID, accountID, analysisDate)
		os.Exit(1)
	}

	// get the report
	lf, err := os.Open(*path)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the requested report: %v\n", err)
		os.Exit(1)
	}
	report := &core.ResourceAccessSummaryReport{}
	err = core.LoadReport(lf, report)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the requested report: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Fprintf(stderr, "Target Analysis: %v, records: %v\n", analysisDate, len(report.Items))
	}

	violations := []ResourceAccessSummary{}
	summaries := BuildResourceAccessSummaries(stderr, report.Items, services, verbose)
	for _, summary := range summaries {
		if !policy.IsCompliant(summary) {
			violations = append(violations, summary)
		}
	}

	views.Display(stdout, stderr, format, violations)
}

type AccessibilityPolicy struct {
	AdminCap  int
	ReadCap   int
	WriteCap  int
	DeleteCap int
}

func (p AccessibilityPolicy) IsCompliant(s ResourceAccessSummary) bool {
	if len(s.PrincipalsByCapability[core.ACCESS_CAPABILITY_RESOURCE_ADMIN]) > p.AdminCap {
		return false
	}
	if len(s.PrincipalsByCapability[core.ACCESS_CAPABILITY_READ_DATA]) > p.ReadCap {
		return false
	}
	if len(s.PrincipalsByCapability[core.ACCESS_CAPABILITY_WRITE_DATA]) > p.WriteCap {
		return false
	}
	if len(s.PrincipalsByCapability[core.ACCESS_CAPABILITY_DELETE_DATA]) > p.DeleteCap {
		return false
	}
	return true
}

type Principal struct {
	ARN  string `json:"principal_arn"`
	Name string `json:"principal_name"`
	Type string `json:"principal_type"`
}

type ResourceAccessSummary struct {
	ServiceName  string `csv:"service_name" json:"service_name"`
	ResourceName string `csv:"resource_name" json:"resource_name"`
	ResourceARN  string `csv:"resource_arn" json:"resource_arn"`

	PrincipalsByCapability map[string][]Principal `csv:"principals_by_capability" json:"principals_by_capability"`
}

func BuildResourceAccessSummaries(stderr io.Writer,
	reportItems []core.ResourceAccessSummaryReportItem,
	services map[string]bool,
	verbose bool) []ResourceAccessSummary {

	indexedSummaries := map[string]ResourceAccessSummary{}
	for _, i := range reportItems {
		var ok bool
		if _, ok = services[i.ServiceName]; !ok {
			if verbose {
				fmt.Fprintf(os.Stderr, "Skipping ReportItem for: %v, %v\n", i.ServiceName, i.ResourceARN)
			}
			continue
		}
		var summary ResourceAccessSummary
		var principal Principal
		var principalList []Principal
		if summary, ok = indexedSummaries[i.ResourceARN]; !ok {
			summary = ResourceAccessSummary{
				ServiceName:            i.ServiceName,
				ResourceName:           i.ResourceName,
				ResourceARN:            i.ResourceARN,
				PrincipalsByCapability: map[string][]Principal{},
			}
		}
		principal = Principal{
			ARN:  i.PrincipalARN,
			Name: i.PrincipalName,
			Type: i.PrincipalType,
		}
		if principalList, ok = summary.PrincipalsByCapability[i.AccessCapability]; !ok {
			principalList = []Principal{}
		}
		principalList = append(principalList, principal)
		summary.PrincipalsByCapability[i.AccessCapability] = principalList
		indexedSummaries[i.ResourceARN] = summary
	}

	summaries := []ResourceAccessSummary{}
	for _, v := range indexedSummaries {
		summaries = append(summaries, v)
	}
	return summaries
}

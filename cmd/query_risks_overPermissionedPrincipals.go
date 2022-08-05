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

// queryRisksOverPermissionedPrincipalsCmd represents the risks command
var queryRisksOverPermissionedPrincipalsCmd = &cobra.Command{
	Use:   "over-permissioned-principals",
	Short: "Show over-permissioned principal risks",
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
		maxRWD, _ := cmd.Flags().GetInt(FLAG_MAX_RWD)

		policy := CapabilityLimitPolicy{
			AdminCap:           maxAdmins,
			ReadWriteDeleteCap: maxRWD}

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

		DoQueryOverPermissionedPrincipals(stdout, stderr,
			reportHome, customerID, accountID, format,
			reportDateTime,
			verbose,
			serviceMap,
			policy)

	},
}

func init() {
	queryRisksCmd.AddCommand(queryRisksOverPermissionedPrincipalsCmd)

	queryRisksOverPermissionedPrincipalsCmd.Flags().StringSlice(FLAG_SERVICE, []string{}, "A list of service names to evaluate")
	queryRisksOverPermissionedPrincipalsCmd.MarkFlagRequired(FLAG_SERVICE)

	queryRisksOverPermissionedPrincipalsCmd.Flags().Int(FLAG_MAX_ADMIN, 5, "The maximum number of resources to which a principal may have administrative access.")
	queryRisksOverPermissionedPrincipalsCmd.Flags().Int(FLAG_MAX_RWD, 5, "The maximum number of resources to which a principal may have read + write + delete access.")
}

func DoQueryOverPermissionedPrincipals(stdout, stderr io.Writer,
	reportHome, customerID, accountID, format string,
	analysisDate *time.Time,
	verbose bool,
	services map[string]bool,
	policy CapabilityLimitPolicy) {

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
	path := db.GetPathForCustomerAccountTimeKind(customerID, accountID, analysisDate, core.REPORT_TYPE_PREFIX_PRINCIPAL_ACCESS_SUMMARIES)
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
	report := &core.PrincipalAccessSummaryReport{}
	err = core.LoadReport(lf, report)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to open the requested report: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Fprintf(stderr, "Target Analysis: %v, records: %v\n", analysisDate, len(report.Items))
	}

	violations := []PrincipalAccessSummary{}
	summaries := BuildPrincipalAccessSummaries(stderr, report.Items, services, verbose)
	for _, summary := range summaries {
		if !policy.IsCompliant(summary) {
			violations = append(violations, summary)
		}
	}

	views.Display(stdout, stderr, format, violations)

}

type CapabilityLimitPolicy struct {
	AdminCap           int
	ReadWriteDeleteCap int
}

func (p CapabilityLimitPolicy) IsCompliant(s PrincipalAccessSummary) bool {
	if len(s.ResourceAccessByCapability[core.ACCESS_CAPABILITY_RESOURCE_ADMIN]) > p.AdminCap {
		return false
	}
	return true
}

type Resource struct {
	ARN         string `json:"resource_arn"`
	ServiceName string `json:"service_name"`
}

type PrincipalAccessSummary struct {
	ARN  string `csv:"principal_arn" json:"principal_arn"`
	Name string `csv:"principal_name" json:"principal_name"`
	Type string `csv:"principal_type" json:"principal_type"`

	ResourceAccessByCapability map[string][]Resource `json:"resources_by_capability"`
}

func BuildPrincipalAccessSummaries(stderr io.Writer, reportItems []core.PrincipalAccessSummaryReportItem, services map[string]bool, verbose bool) []PrincipalAccessSummary {
	indexedSummaries := map[string]PrincipalAccessSummary{}
	for _, i := range reportItems {
		var ok bool
		if _, ok = services[i.ServiceName]; !ok {
			if verbose {
				fmt.Fprintf(stderr, "Skipping ReportItem for: %v, %v\n", i.ServiceName, i.PrincipalARN)
			}
			continue
		}

		var summary PrincipalAccessSummary
		var resource Resource
		var resourceList []Resource
		if summary, ok = indexedSummaries[i.PrincipalARN]; !ok {
			summary = PrincipalAccessSummary{
				ARN:                        i.PrincipalARN,
				Name:                       i.PrincipalName,
				Type:                       i.PrincipalType,
				ResourceAccessByCapability: map[string][]Resource{},
			}
		}
		resource = Resource{
			ARN:         i.ResourceARN,
			ServiceName: i.ServiceName,
		}
		if resourceList, ok = summary.ResourceAccessByCapability[i.AccessCapability]; !ok {
			resourceList = []Resource{}
		}
		resourceList = append(resourceList, resource)
		summary.ResourceAccessByCapability[i.AccessCapability] = resourceList
		indexedSummaries[i.PrincipalARN] = summary
	}

	summaries := []PrincipalAccessSummary{}
	for _, v := range indexedSummaries {
		summaries = append(summaries, v)
	}
	return summaries
}

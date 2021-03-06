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

package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type ReportSet struct {
	CustomerID string
	Account    string

	IndexedByMonth map[string]map[string][]Report
	Index          map[MonthKey][]Report
	Set            []Report
}

// MonthKey is a small structed used to structure the keyspace for ReportSetIndexes.
type MonthKey struct {
	Year, Month string
}
type ReportSetIndex map[MonthKey][]Report

func BuildIndex(set ReportSet) ReportSetIndex {
	index := ReportSetIndex{}
	for _, r := range set.Set {
		k := MonthKey{strconv.Itoa(r.Timestamp.Year()), r.Timestamp.Format(MONTH_TIMESTAMP_LAYOUT)}
		if _, ok := index[k]; !ok {
			index[k] = append(index[k], r)
		}
	}
	return index
}

// Reports represents a collection of reports generated for the
// same customer, account, and at the same reporting time.
// This design assumes that all reports related to the same run
// will have the same minute field in the file name.
type Report struct {
	Bucket     string
	CustomerID string
	Account    string
	Timestamp  time.Time
}

func (r Report) PrincipalAccessS3ObjectKey() string {
	return r.reportS3ObjectKey(`principal-access-summaries`)
}

func (r Report) ResourceAccessS3ObjectKey() string {
	return r.reportS3ObjectKey(`resource-access-summaries`)
}

func (r Report) PrincipalsS3ObjectKey() string {
	return r.reportS3ObjectKey(`principals`)
}

func (r Report) ResourcesS3ObjectKey() string {
	return r.reportS3ObjectKey(`resources`)
}

func (r Report) ResourceAccessAuditS3ObjectKey() string {
	return fmt.Sprintf(
		REPORT_LOCATION_XLSX_FQ_PATTERN,
		r.CustomerID,
		r.Account,
		strconv.Itoa(r.Timestamp.Year()),
		r.Timestamp.Format(MONTH_TIMESTAMP_LAYOUT),
		`resource-access-audit`,
		r.Timestamp.Format(FILENAME_TIMESTAMP_LAYOUT))
}

func (r Report) reportS3ObjectKey(name string) string {
	return fmt.Sprintf(
		REPORT_LOCATION_CSV_FQ_PATTERN,
		r.CustomerID,
		r.Account,
		strconv.Itoa(r.Timestamp.Year()),
		r.Timestamp.Format(MONTH_TIMESTAMP_LAYOUT),
		name,
		r.Timestamp.Format(FILENAME_TIMESTAMP_LAYOUT))
}

type ResourcesReportItem struct {
	AnalysisTime time.Time `csv:"analysis_time" json:"analysis_time"`
	ResourceName string    `csv:"resource_name" json:"resource_name"`
	ResourceARN  string    `csv:"resource_arn" json:"resource_arn"`
	ResourceType string    `csv:"resource_type" json:"resource_type"`

	ResourceTagBusinessUnit    string `csv:"resource_tag_business_unit" json:"resource_tag_business_unit"`
	ResourceTagEnvironment     string `csv:"resource_tag_environment" json:"resource_tag_environment"`
	ResourceTagOwner           string `csv:"resource_tag_owner" json:"resource_tag_owner"`
	ResourceTagConfidentiality string `csv:"resource_tag_confidentiality" json:"resource_tag_confidentiality"`
	ResourceTagIntegrity       string `csv:"resource_tag_integrity" json:"resource_tag_integrity"`
	ResourceTagAvailability    string `csv:"resource_tag_availability" json:"resource_tag_availability"`
	ResourceTags               string `csv:"resource_tags" json:"resource_tags"`
}

func (i ResourcesReportItem) Equivalent(t ResourcesReportItem) bool {
	if i.ResourceName != t.ResourceName ||
		i.ResourceARN != t.ResourceARN ||
		i.ResourceType != t.ResourceType ||
		i.ResourceTagBusinessUnit != t.ResourceTagBusinessUnit ||
		i.ResourceTagEnvironment != t.ResourceTagEnvironment ||
		i.ResourceTagOwner != t.ResourceTagOwner ||
		i.ResourceTagConfidentiality != t.ResourceTagConfidentiality ||
		i.ResourceTagIntegrity != t.ResourceTagIntegrity ||
		i.ResourceTagAvailability != t.ResourceTagAvailability ||
		i.ResourceTags != t.ResourceTags {
		return false
	}
	return true
}

func UnmarshalResourcesReportItem(in []string) (o ResourcesReportItem, err error) {
	if len(in) != 11 {
		err = fmt.Errorf(`invalid Resources Report Item record length`)
		return
	}
	o.AnalysisTime, err = time.Parse(time.RFC3339Nano, in[0])
	if err != nil {
		return
	}
	o.ResourceName = in[1]
	o.ResourceARN = in[2]
	o.ResourceType = in[3]
	o.ResourceTagBusinessUnit = in[4]
	o.ResourceTagEnvironment = in[5]
	o.ResourceTagOwner = in[6]
	o.ResourceTagConfidentiality = in[7]
	o.ResourceTagIntegrity = in[8]
	o.ResourceTagAvailability = in[9]
	o.ResourceTags = in[10]
	return
}

type PrincipalsReportItem struct {
	AnalysisTime        time.Time `csv:"analysis_time" json:"analysis_time"`
	PrincipalName       string    `csv:"principal_name" json:"principal_name"`
	PrincipalARN        string    `csv:"principal_arn" json:"principal_arn"`
	PrincipalType       string    `csv:"principal_type" json:"principal_type"`
	PrincipalIsIAMAdmin bool      `csv:"principal_is_iam_admin" json:"principal_is_iam_admin"`
	PrincipalLastUsed   string    `csv:"principal_last_used" json:"principal_last_used"`

	PrincipalTagBusinessUnit string `csv:"principal_tag_business_unit" json:"principal_tag_business_unit"`
	PrincipalTagEnvironment  string `csv:"principal_tag_environment" json:"principal_tag_environment"`
	PrincipalTagUsedBy       string `csv:"principal_tag_used_by" json:"principal_tag_used_by"`
	PrincipalTags            string `csv:"principal_tags" json:"principal_tags"`

	PasswordLastUsed    string `csv:"password_last_used" json:"password_last_used"`
	PasswordLastRotated string `csv:"password_last_rotated" json:"password_last_rotated"`
	PasswordState       string `csv:"password_state" json:"password_state"`

	AccessKey1LastUsed    string `csv:"access_key_1_last_used" json:"access_key_1_last_used"`
	AccessKey1LastRotated string `csv:"access_key_1_last_rotated" json:"access_key_1_last_rotated"`
	AccessKey1State       string `csv:"access_key_1_state" json:"access_key_1_state"`
	AccessKey2LastUsed    string `csv:"access_key_2_last_used" json:"access_key_2_last_used"`
	AccessKey2LastRotated string `csv:"access_key_2_last_rotated" json:"access_key_2_last_rotated"`
	AccessKey2State       string `csv:"access_key_2_state" json:"access_key_2_state"`
}

func (i PrincipalsReportItem) Equivalent(t PrincipalsReportItem) bool {
	if i.PrincipalName != t.PrincipalName ||
		i.PrincipalARN != t.PrincipalARN ||
		i.PrincipalType != t.PrincipalType ||
		i.PrincipalIsIAMAdmin != t.PrincipalIsIAMAdmin ||
		i.PrincipalLastUsed != t.PrincipalLastUsed ||
		i.PrincipalTagBusinessUnit != t.PrincipalTagBusinessUnit ||
		i.PrincipalTagEnvironment != t.PrincipalTagEnvironment ||
		i.PrincipalTagUsedBy != t.PrincipalTagUsedBy ||
		i.PrincipalTags != t.PrincipalTags ||
		i.PasswordLastUsed != t.PasswordLastUsed ||
		i.PasswordLastRotated != t.PasswordLastRotated ||
		i.PasswordState != t.PasswordState ||
		i.AccessKey1LastUsed != t.AccessKey1LastUsed ||
		i.AccessKey1LastRotated != t.AccessKey1LastRotated ||
		i.AccessKey1State != t.AccessKey1State ||
		i.AccessKey2LastUsed != t.AccessKey2LastUsed ||
		i.AccessKey2LastRotated != t.AccessKey2LastRotated ||
		i.AccessKey2State != t.AccessKey2State {
		return false
	}
	return true
}

func UnmarshalPrincipalsReportItem(in []string) (o PrincipalsReportItem, err error) {
	if len(in) != 19 {
		err = &IllegalArgumentError{`in`, `invalid PrincipalsReportItem entry`}
		return
	}
	o.AnalysisTime, err = time.Parse(time.RFC3339Nano, in[0])
	if err != nil {
		return
	}
	o.PrincipalName = in[1]
	o.PrincipalARN = in[2]
	o.PrincipalType = in[3]
	o.PrincipalIsIAMAdmin, _ = strconv.ParseBool(in[4])
	o.PrincipalLastUsed = in[5]
	o.PrincipalTagBusinessUnit = in[6]
	o.PrincipalTagEnvironment = in[7]
	o.PrincipalTagUsedBy = in[8]
	o.PrincipalTags = in[9]
	o.PasswordLastUsed = in[10]
	o.PasswordLastRotated = in[11]
	o.PasswordState = in[12]
	o.AccessKey1LastUsed = in[13]
	o.AccessKey1LastRotated = in[14]
	o.AccessKey1State = in[15]
	o.AccessKey2LastUsed = in[16]
	o.AccessKey2LastRotated = in[17]
	o.AccessKey2State = in[18]
	return
}

type PrincipalAccessSummaryReportItem struct {
	AnalysisTime     time.Time `csv:"analysis_time" json:"analysis_time"`
	PrincipalName    string    `csv:"principal_name" json:"principal_name"`
	PrincipalARN     string    `csv:"principal_arn" json:"principal_arn"`
	PrincipalType    string    `csv:"principal_type" json:"principal_type"`
	PrincipalTags    string    `csv:"principal_tags" json:"principal_tags"`
	ServiceName      string    `csv:"service_name" json:"service_name"`
	AccessCapability string    `csv:"access_capability" json:"access_capability"`
	ResourceARN      string    `csv:"resource_arn" json:"resource_arn"`
}

func (i PrincipalAccessSummaryReportItem) Equivalent(t PrincipalAccessSummaryReportItem) bool {
	if i.PrincipalName != t.PrincipalName ||
		i.PrincipalARN != t.PrincipalARN ||
		i.PrincipalType != t.PrincipalType ||
		i.PrincipalTags != t.PrincipalTags ||
		i.ServiceName != t.ServiceName ||
		i.AccessCapability != t.AccessCapability ||
		i.ResourceARN != t.ResourceARN {
		return false
	}
	return true
}

func UnmarshalPrincipalAccessSummaryReportItem(in []string) (o PrincipalAccessSummaryReportItem, err error) {
	if len(in) != 8 {
		err = &IllegalArgumentError{`in`, `invalid PrincipalAccessReportItem entry`}
		return
	}
	o.AnalysisTime, err = time.Parse(time.RFC3339Nano, in[0])
	if err != nil {
		return
	}
	o.PrincipalName = in[1]
	o.PrincipalARN = in[2]
	o.PrincipalType = in[3]
	o.PrincipalTags = in[4]
	o.ServiceName = in[5]
	o.AccessCapability = in[6]
	o.ResourceARN = in[7]
	return
}

type ResourceAccessSummaryReportItem struct {
	AnalysisTime     time.Time `csv:"analysis_time" json:"analysis_time"`
	ServiceName      string    `csv:"service_name" json:"service_name"`
	ResourceName     string    `csv:"resource_name" json:"resource_name"`
	ResourceARN      string    `csv:"resource_arn" json:"resource_arn"`
	AccessCapability string    `csv:"access_capability" json:"access_capability"`
	PrincipalType    string    `csv:"principal_type" json:"principal_type"`
	PrincipalName    string    `csv:"principal_name" json:"principal_name"`
	PrincipalARN     string    `csv:"principal_arn" json:"principal_arn"`

	ResourceTagConfidentiality string `csv:"resource_tag_confidentiality" json:"resource_tag_confidentiality"`
}

func (i ResourceAccessSummaryReportItem) Equivalent(t ResourceAccessSummaryReportItem) bool {
	if i.PrincipalName != t.PrincipalName ||
		i.PrincipalARN != t.PrincipalARN ||
		i.PrincipalType != t.PrincipalType ||
		i.ServiceName != t.ServiceName ||
		i.AccessCapability != t.AccessCapability ||
		i.ResourceName != t.ResourceName ||
		i.ResourceARN != t.ResourceARN ||
		i.ResourceTagConfidentiality != t.ResourceTagConfidentiality {
		return false
	}
	return true
}

func UnmarshalResourceAccessSummaryReportItem(in []string) (o ResourceAccessSummaryReportItem, err error) {
	if len(in) != 9 {
		err = &IllegalArgumentError{`in`, `invalid ResourceAccessReportItem entry`}
		return
	}
	o.AnalysisTime, err = time.Parse(time.RFC3339Nano, in[0])
	if err != nil {
		return
	}
	o.ServiceName = in[1]
	o.ResourceName = in[2]
	o.ResourceARN = in[3]
	o.AccessCapability = in[4]
	o.PrincipalType = in[5]
	o.PrincipalName = in[6]
	o.PrincipalARN = in[7]
	o.ResourceTagConfidentiality = in[8]
	return
}

// LoadReport reads all records from the provided Reader as CSV and aggregates those records using
// the provided Collector.
func LoadReport(in io.Reader, c Collector) error {
	rr := csv.NewReader(in)
	records, err := rr.ReadAll()
	if err != nil {
		return err
	}

	for i, v := range records {
		// skip the header row
		if i == 0 {
			continue
		}
		if err := c.Collect(v); err != nil {
			return err
		}
	}
	return nil
}

// Collector describes record-aggregating recievers. A Collector implementation should collect a
// specific type of record. For example a ResourceAccessSummaryReport is a Collector that will
// attempt to parse a ResourceAccessSummaryReportItem from the provided string slice and append
// that record to the report's internal aggregation.
type Collector interface {
	Collect(in []string) error
}

// ResourceAccessSummaryReport is a ResourceAccessSummaryReportItem collector.
type ResourceAccessSummaryReport struct {
	Items []ResourceAccessSummaryReportItem
}

// Collect will attempt to parse a ResourceAccessSummaryReportItem and append it to the
// ResourceAccessSummaryReport internal aggregation.
func (r *ResourceAccessSummaryReport) Collect(in []string) error {
	if r.Items == nil {
		r.Items = []ResourceAccessSummaryReportItem{}
	}
	ri, err := UnmarshalResourceAccessSummaryReportItem(in)
	if err != nil {
		return err
	}
	r.Items = append(r.Items, ri)
	return nil
}

// PrincipalAccessSummaryReport is a PrincipalAccessSummaryReportItem collector.
type PrincipalAccessSummaryReport struct {
	Items []PrincipalAccessSummaryReportItem
}

// Collect will attempt to parse a PrincipalAccessSummaryReportItem and append it to the
// PrincipalAccessSummaryReport internal aggregation.
func (r *PrincipalAccessSummaryReport) Collect(in []string) error {
	if r.Items == nil {
		r.Items = []PrincipalAccessSummaryReportItem{}
	}
	ri, err := UnmarshalPrincipalAccessSummaryReportItem(in)
	if err != nil {
		return err
	}
	r.Items = append(r.Items, ri)
	return nil
}

// PrincipalReport is a PrincipalReportItem collector.
type PrincipalsReport struct {
	Items []PrincipalsReportItem
}

// Collect will attempt to parse a PrincipalReportItem and append it to the
// PrincipalReport internal aggregation.
func (r *PrincipalsReport) Collect(in []string) error {
	if r.Items == nil {
		r.Items = []PrincipalsReportItem{}
	}
	ri, err := UnmarshalPrincipalsReportItem(in)
	if err != nil {
		return err
	}
	r.Items = append(r.Items, ri)
	return nil
}

// ResourceReport is a ResourceReportItem collector.
type ResourcesReport struct {
	Items []ResourcesReportItem
}

// Collect will attempt to parse a ResourceReportItem and append it to the
// ResourceReport internal aggregation.
func (r *ResourcesReport) Collect(in []string) error {
	if r.Items == nil {
		r.Items = []ResourcesReportItem{}
	}
	ri, err := UnmarshalResourcesReportItem(in)
	if err != nil {
		return err
	}
	r.Items = append(r.Items, ri)
	return nil
}

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

var (
	REPORT_LOCATION_PREFIX           = `customers/`
	REPORT_LOCATION_DELIMITER        = `/`
	REPORT_LOCATION_CSV_FQ_PATTERN   = `customers/%s/reports/aws/%s/%s/%s/%s.%s.csv`
	REPORT_LOCATION_XLSX_FQ_PATTERN  = `customers/%s/reports/aws/%s/%s/%s/%s.%s.xlsx`
	REPORT_LOCATION_CUSTOMER_PATTERN = `customers/%s/reports/aws/`
	REPORT_LOCATION_ACCOUNT_PATTERN  = `customers/%s/reports/aws/%s/`
	REPORT_LOCATION_MONTH_PATTERN    = `customers/%s/reports/aws/%s/%s/%s`
)

const (
	EXT_CSV  = `csv`
	EXT_XLSX = `xlsx`
)

// report file name prefixes
const (
	REPORT_TYPE_PREFIX_PRINCIPALS                 = `principals`
	REPORT_TYPE_PREFIX_RESOURCES                  = `resources`
	REPORT_TYPE_PREFIX_PRINCIPAL_ACCESS_SUMMARIES = `principal-access-summaries`
	REPORT_TYPE_PREFIX_RESOURCE_ACCESS_SUMMARIES  = `resource-access-summaries`
)

const (
	_ = iota
	FILENAME_POSITION_CID
	_
	_
	FILENAME_POSITION_ACCOUNT
	FILENAME_POSITION_YEAR
	FILENAME_POSITION_MONTH
	FILENAME_POSITION_FILE
)
const (
	FILENAME_TIMESTAMP_ANALYSIS_DATE_LAYOUT = `2006-01-02`
	FILENAME_TIMESTAMP_LAYOUT               = "2006-01-02-1504"
	MONTH_TIMESTAMP_LAYOUT                  = "01"
	LATEST                                  = "latest"
)

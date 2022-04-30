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
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func List(o io.Writer, cfg aws.Config, bucket, customerID, account string) error {
	if len(customerID) <= 0 {
		// no customers specified, list the customers
		return listCustomers(o, cfg, bucket)
	} else if len(account) <= 0 {
		// no account specified, list accounts
		return listAccounts(o, cfg, bucket, customerID)
	} else {
		// list objects matching some pattern
		if reports, err := listObjects(cfg, bucket, customerID, account); err != nil {
			return err
		} else {
			return displayReports(o, reports)
		}
	}
}

func listCustomers(o io.Writer, cfg aws.Config, bucket string) error {
	client := s3.NewFromConfig(cfg)
	pages := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Delimiter: &REPORT_LOCATION_DELIMITER,
		Bucket:    &bucket,
		Prefix:    &REPORT_LOCATION_PREFIX})

	for pages.HasMorePages() {
		resp, err := pages.NextPage(context.TODO())
		if err != nil {
			return err
		}

		for _, p := range resp.CommonPrefixes {
			s := strings.Split(*p.Prefix, REPORT_LOCATION_DELIMITER)
			if len(s) < 2 {
				// malformed prefix
			}
			fmt.Fprintln(o, s[1])
		}
	}
	return nil
}

func listAccounts(o io.Writer, cfg aws.Config, bucket, customerID string) error {
	prefix := fmt.Sprintf(REPORT_LOCATION_CUSTOMER_PATTERN, customerID)
	client := s3.NewFromConfig(cfg)
	pages := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Delimiter: &REPORT_LOCATION_DELIMITER,
		Bucket:    &bucket,
		Prefix:    &prefix})

	for pages.HasMorePages() {
		resp, err := pages.NextPage(context.TODO())
		if err != nil {
			return err
		}

		for _, p := range resp.CommonPrefixes {
			s := strings.Split(*p.Prefix, REPORT_LOCATION_DELIMITER)
			if len(s) < 5 {
				// malformed prefix
			}
			fmt.Fprintln(o, s[4])
		}
	}
	return nil
}

func listObjects(cfg aws.Config, bucket, customerID, account string) (ReportSet, error) {
	prefix := fmt.Sprintf(REPORT_LOCATION_ACCOUNT_PATTERN, customerID, account)

	reports := ReportSet{CustomerID: customerID, Account: account}
	index := map[time.Time]Report{}

	client := s3.NewFromConfig(cfg)
	pages := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix})

	for pages.HasMorePages() {
		resp, err := pages.NextPage(context.TODO())
		if err != nil {
			return reports, err
		}

		for _, o := range resp.Contents {
			if !strings.HasSuffix(*o.Key, `csv`) {
				continue
			}

			s := strings.Split(*o.Key, REPORT_LOCATION_DELIMITER)
			if len(s) < 7 {
				// TODO malformed prefix
			}

			// parse the filename
			fparts := strings.Split(s[FILENAME_POSITION_FILE], ".")
			if len(fparts) < 3 {
				// TODO malformed filename
			}
			rts, err := time.Parse(FILENAME_TIMESTAMP_LAYOUT, fparts[1])
			if err != nil {
				// TODO malformed report filename
			}

			if _, ok := index[rts]; !ok {
				fresh := Report{
					Bucket:     bucket,
					CustomerID: customerID,
					Account:    account,
					Timestamp:  rts}
				reports.Set = append(reports.Set, fresh)
				index[rts] = fresh
			}
		}
	}

	return reports, nil
}

func displayReports(o io.Writer, reports ReportSet) error {
	sort.Slice(reports.Set, func(p, q int) bool {
		return reports.Set[p].Timestamp.Before(reports.Set[q].Timestamp)
	})
	for _, v := range reports.Set {
		fmt.Fprintf(o, "%s\n", v.Timestamp.Format(FILENAME_TIMESTAMP_LAYOUT))
	}
	return nil
}

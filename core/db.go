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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DB struct {
	Customers map[string]Customer
}

func (db *DB) Dump(o io.Writer, isSummary bool) {
	for _, c := range db.Customers {
		fmt.Fprintf(o, "%v #%v\n", c.CustomerID, len(c.Accounts))
		for _, a := range c.Accounts {
			fmt.Fprintf(o, "\t%v #%v\n", a.AccountID, len(a.Reports))
			if !isSummary {
				for t, r := range a.Reports {
					fmt.Fprintf(o, "\t\t%v: %v\n", t, r)
				}
			}
		}
	}
}

func (db *DB) Size() (total int) {
	for _, c := range db.Customers {
		for _, a := range c.Accounts {
			total += len(a.Reports)
		}
	}
	return
}

func (db *DB) Sizes() (total int, accounts int, customers int) {
	customers = len(db.Customers)
	for _, c := range db.Customers {
		accounts += len(c.Accounts)
		for _, a := range c.Accounts {
			total += len(a.Reports)
		}
	}
	return
}

func (db *DB) GetPathForCustomerAccountTimeKind(customerID, accountID string, ts time.Time, kind string) *string {
	var (
		customer Customer
		account  Account
		report   LocalReport
		path     string
		ok       bool
	)
	if customer, ok = db.Customers[customerID]; !ok {
		return nil
	}
	if account, ok = customer.Accounts[accountID]; !ok {
		return nil
	}
	if report, ok = account.Reports[ts.Truncate(24*time.Hour)]; !ok {
		return nil
	}
	if path, ok = report.pathByKind[kind]; !ok {
		return nil
	}
	return &path
}

func (db *DB) GetPathForCustomerAccountLatestKind(customerID, accountID, kind string) *string {
	var (
		customer Customer
		account  Account
		report   LocalReport
		path     string
		ok       bool
	)
	if customer, ok = db.Customers[customerID]; !ok {
		return nil
	}
	if account, ok = customer.Accounts[accountID]; !ok {
		return nil
	}
	report = account.Latest()
	if path, ok = report.pathByKind[kind]; !ok {
		return nil
	}
	return &path
}

func (db *DB) AllPaths() []string {
	out := []string{}
	for _, c := range db.Customers {
		for _, a := range c.Accounts {
			for _, r := range a.Reports {
				for _, p := range r.pathByKind {
					out = append(out, p)
				}
			}
		}
	}
	return out
}

func (db *DB) AllPathsByCustomerAccount(customerID, accountID string) []string {
	out := []string{}
	for _, c := range db.Customers {
		if c.CustomerID != customerID {
			continue
		}
		for _, a := range c.Accounts {
			if a.AccountID != accountID {
				continue
			}
			for _, r := range a.Reports {
				for _, p := range r.pathByKind {
					out = append(out, p)
				}
			}
		}
	}
	return out
}

type Customer struct {
	CustomerID string
	Accounts   map[string]Account
}

type Account struct {
	AccountID string
	Reports   map[time.Time]LocalReport
}

func (a *Account) Latest() LocalReport {
	lt := time.UnixMicro(0)
	var latest LocalReport
	for t, r := range a.Reports {
		if lt.Before(t) {
			lt = t
			latest = r
		}
	}
	return latest
}

type LocalReport struct {
	CustomerID string
	Account    string
	Timestamp  time.Time
	pathByKind map[string]string
}

func LoadLocalDB(root string) (DB, error) {
	out := DB{Customers: map[string]Customer{}}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		return dbDirWalker(&out, root, path, info, err)
	})
	return out, err
}

type ReportTypeSelector []string

// LoadS3DB enumerates and pulls metadata for all customers, accounts, and reports in
// the specified S3 bucket. It does however, skip unknown report types.
func LoadS3DB(client s3.ListObjectsV2APIClient, bucket string, selector ReportTypeSelector) (DB, error) {
	out := DB{Customers: map[string]Customer{}}
	pages := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{Bucket: &bucket})
	for pages.HasMorePages() {
		resp, err := pages.NextPage(context.TODO())
		if err != nil {
			return out, err
		}

		for _, v := range resp.Contents {
			isSelected := false
			for _, t := range selector {
				if !isSelected && strings.HasSuffix(*v.Key, t) {
					isSelected = true
				}
			}
			if !isSelected {
				continue
			}

			// there is some disagreement about if this should have 7 or 8 parts in S3
			parts := strings.Split(*v.Key, REPORT_LOCATION_DELIMITER)
			if len(parts) != 8 {
				continue
			}

			var ok bool
			var customer Customer
			if customer, ok = out.Customers[parts[DB_INDEX_POSITION_CUSTOMERID]]; !ok {
				customer = Customer{CustomerID: parts[DB_INDEX_POSITION_CUSTOMERID], Accounts: map[string]Account{}}
				out.Customers[customer.CustomerID] = customer
			}
			// retrieve / initialize the account entry
			var account Account
			if account, ok = customer.Accounts[parts[DB_INDEX_POSITION_ACCOUNT]]; !ok {
				account = Account{
					AccountID: parts[DB_INDEX_POSITION_ACCOUNT],
					Reports:   map[time.Time]LocalReport{}}
				customer.Accounts[account.AccountID] = account
			}

			// parse out the type and date of the individual report file
			base := parts[DB_INDEX_POSITION_FILE]
			baseParts := strings.Split(base, `.`)
			if len(baseParts) != 3 {
				// return fmt.Errorf(`invalid report filename, invalid filename structure, %v`, base)
				continue
			}
			if baseParts[1] == LATEST {
				continue
			}
			reportTime, err := time.Parse(FILENAME_TIMESTAMP_LAYOUT, baseParts[1])
			if err != nil {
				// return fmt.Errorf(`invalid report filename, invalid timestamp`)
				continue
			}
			reportTimeTruncated := reportTime.Truncate(24 * time.Hour)
			var report LocalReport
			if report, ok = account.Reports[reportTimeTruncated]; !ok {
				report = LocalReport{
					CustomerID: customer.CustomerID,
					Account:    account.AccountID,
					Timestamp:  reportTimeTruncated,
					pathByKind: map[string]string{}}
				account.Reports[reportTimeTruncated] = report
			}
			report.pathByKind[baseParts[0]] = *v.Key
		}
	}
	return out, nil
}

// dbDirWalker is to be used with a wrapper for filepath.Walk and is to be invoked
// for each file under some specific point in the file tree. This function adds
// customer, account, and report records to a provided DB instance. Access to the
// provided DB instance is not synchronized. For that reason this func should not
// be called in a goroutine.
func dbDirWalker(out *DB, root, path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	// skip directories, this routine parses the filename anyway
	if info.IsDir() {
		return nil
	}
	// parse the path
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}
	parts := strings.Split(rel, string(os.PathSeparator))
	if len(parts) != 8 {
		return nil
	}

	var ok bool
	// retrieve / initialize the customer entry
	var customer Customer
	if customer, ok = out.Customers[parts[DB_INDEX_POSITION_CUSTOMERID]]; !ok {
		customer = Customer{CustomerID: parts[DB_INDEX_POSITION_CUSTOMERID], Accounts: map[string]Account{}}
		out.Customers[customer.CustomerID] = customer
	}

	// retrieve / initialize the account entry
	var account Account
	if account, ok = customer.Accounts[parts[DB_INDEX_POSITION_ACCOUNT]]; !ok {
		account = Account{
			AccountID: parts[DB_INDEX_POSITION_ACCOUNT],
			Reports:   map[time.Time]LocalReport{}}
		customer.Accounts[account.AccountID] = account
	}

	// parse out the type and date of the individual report file
	base := parts[DB_INDEX_POSITION_FILE]
	baseParts := strings.Split(base, `.`)
	if len(baseParts) != 3 {
		return fmt.Errorf(`invalid report filename, invalid filename structure`)
	}
	if baseParts[1] == LATEST {
		return nil
	}
	reportTime, err := time.Parse(FILENAME_TIMESTAMP_LAYOUT, baseParts[1])
	if err != nil {
		return fmt.Errorf(`invalid report filename, invalid timestamp`)
	}
	reportTimeTruncated := reportTime.Truncate(24 * time.Hour)
	var report LocalReport
	if report, ok = account.Reports[reportTimeTruncated]; !ok {
		report = LocalReport{
			CustomerID: customer.CustomerID,
			Account:    account.AccountID,
			Timestamp:  reportTimeTruncated,
			pathByKind: map[string]string{}}
		account.Reports[reportTimeTruncated] = report
	}
	report.pathByKind[baseParts[0]] = rel
	return nil
}

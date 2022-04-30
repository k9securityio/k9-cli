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

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Sync(stdout, stderr io.Writer,
	remote DB,
	downloader *manager.Downloader,
	bucket, customerID, account string,
	concurrency int,
	dryrun, verbose bool) error {

	// TODO input validation

	// collector slice for errors that occur in processing
	errs := []error{}

	// setup concurrent downloading harness
	semaphore := make(chan int, concurrency) // buffered channel
	download := func(key string, w WriterAtCloser) {
		defer func() {
			w.Close()
			<-semaphore
		}()
		if !dryrun {
			if _, err := downloader.Download(
				context.TODO(), w, &s3.GetObjectInput{Bucket: &bucket, Key: &key}); err != nil {
				// Don't do this... append is not synchronized
				// errs = append(errs, err)
				// TODO use a channel to aggregate the errors
				fmt.Fprintln(stderr, err)
			}
		}
		if verbose {
			fmt.Fprintln(stderr, key)
		}
	}

	// get the target payloads
	pathsToSync := remote.AllPathsByCustomerAccount(customerID, account)

	for _, r := range pathsToSync {
		semaphore <- 1
		folder, _ := filepath.Split(r)
		if err := os.MkdirAll(folder, 0750); err != nil {
			errs = append(errs, err)
			continue
		}
		f, err := os.Create(r)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		go download(r, f)
	}

	// TODO if any errors during retrival or writing aggregate those into an aggregate error
	if len(errs) > 0 {
		return &AggregateError{true, errs}
	}
	return nil
}

type WriterAtCloser interface {
	io.WriterAt
	io.Closer
}

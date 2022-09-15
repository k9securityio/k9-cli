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
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type analyzeRequestBody struct {
	CustomerID string `json:"customerId"`
	Account    string `json:"accountId"`
}

// Example:
// {"customerId": "C10000", "accountId": "139710491120", "executionId": "ondemand-C10000-139710491120-2022-09-15_TV49"}
type analyzeResponseBody struct {
	CustomerID string `json:"customerId"`
	Account    string `json:"accountId"`
	ExecutionID string `json:"executionId"`
}

func AnalyzeAccount(o io.Writer, cfg aws.Config, apiHost, customerID, account string) error {
	url := fmt.Sprintf("https://%s/analysis/account", apiHost)

	requestBody := analyzeRequestBody{
		CustomerID: customerID,
		Account:    account,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not build API request: %s\n", err)
		return err
	}
	now := time.Now()
	request.Header.Set("Date", now.Format(time.RFC3339))
	request.Header.Set("Content-Type", "application/json")

	err = signApiRequest(cfg, request, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not sign API request: %s\n", err)
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not execute API request: %s\n", err)
		return err
	} else {
		defer response.Body.Close()
	}

	body, _ := ioutil.ReadAll(response.Body)
	responseBodyStr := string(body)
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusAccepted {
		analyzeResponse := analyzeResponseBody{}
		err = json.Unmarshal(body, &analyzeResponse)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not deserialize API response: %s\n", responseBodyStr)
			return err
		}

		fmt.Fprintf(o, "Started analysis for %s account %s with execution ID: %s\n",
			customerID,
			account,
			analyzeResponse.ExecutionID)
	} else {
		//fmt.Println("response Headers:", response.Header)
		fmt.Fprintf(os.Stderr,"Analyze API Response Status: %s\n", response.Status)
		fmt.Fprintf(os.Stderr, "Could not start analysis for %s account %s.  API Response: %s\n",
			customerID,
			account,
			responseBodyStr)
	}
	return nil
}

// Sign a request to the k9 Security AWS API gateway endpoint with an AWS v4 signature
// using the aws-sdk-go-v2 v4.SignHTTP function.
//
// The request should be fully-built and ready to send to the gateway.
//
// The provided request will be modified in place.
func signApiRequest(cfg aws.Config, request *http.Request, signingTime time.Time) error {
	ctx := context.TODO()
	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}

	hash, err := createPayloadHash(request)
	if err != nil {
		return err
	}

	signer := v4.NewSigner()
	err = signer.SignHTTP(ctx, credentials, request, hash, "execute-api", cfg.Region, signingTime)

	if err != nil {
		return err
	}

	// fmt.Printf("Authorization: %s\n", request.Header.Get("Authorization"))
	return err
}

// Create a hex-encoded SHA256 hash of the request payload.  Use when calculating an AWS v4 signature.
func createPayloadHash(req *http.Request) (string, error) {
	// from https://github.com/yuizho/salon/blob/a159bcfeb263cb6502403f682c138286a2d7bb1f/backend/lambda/mutate-user/appsync/client.go
	body, err := req.GetBody()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(body)

	if err != nil {
		return "", err
	}
	b := sha256.Sum256(buf.Bytes())
	return hex.EncodeToString(b[:]), nil
}

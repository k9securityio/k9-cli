package core

import (
	"testing"
	"time"
)

func TestExtractReportTime(t *testing.T) {
	//customers/C10001/reports/aws/139710491120/2021/05/resources.2021-05-30-0750.csv
	//customers/C10001/reports/aws/139710491120/2021/05/principals.2021-05-30-0750.csv
	//customers/C10001/reports/aws/139710491120/2021/06/principals.2021-06-08-0755.csv
	//customers/C10001/reports/aws/139710491120/2021/05/resource-access-summaries.2021-05-30-0750.csv
	//customers/C10001/reports/aws/139710491120/2021/05/principal-access-summaries.2021-05-30-0750.csv
	reportTime, err := extractReportTimeFromKey("customers/C10001/reports/aws/139710491120/2021/05/resources.2021-05-30-0750.csv")
	if err != nil {
		t.Fatalf(`Expected err to be nil, was %v`, err)
	}

	expectTime, _ := time.Parse("2006-01-02-1504", "2021-05-30-0751")

	if expectTime != reportTime {
		t.Fatalf(`Expected report time to be %v, but was %v`, expectTime, reportTime)
	}
}

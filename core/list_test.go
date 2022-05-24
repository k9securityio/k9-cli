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
	cases := map[string]struct {
		Key          string
		ExpectedTime time.Time
		ExpectedErr  bool
	}{
		`valid resources file - fully qualified dt`: {
			Key:          `customers/C10001/reports/aws/139710491120/2021/05/resources.2021-05-30-0750.csv`,
			ExpectedTime: parseTime("2021-05-30-0750"),
			ExpectedErr:  false,
		},
		`valid principals file - fully qualified dt`: {
			Key:          `customers/C10001/reports/aws/139710491120/2021/06/principals.2021-06-11-1755.csv`,
			ExpectedTime: parseTime("2021-06-11-1755"),
			ExpectedErr:  false,
		},
		`valid principal access summary file - fully qualified dt`: {
			Key:          `customers/C10001/reports/aws/139710491120/2021/11/principal-access-summaries.2021-11-30-2350.csv`,
			ExpectedTime: parseTime("2021-11-30-2350"),
			ExpectedErr:  false,
		},
	}
	for l, c := range cases {
		actualTime, err := extractReportTimeFromKey(c.Key)
		if err != nil {
			t.Fatalf(`ExpectedTime err to be nil, was %v`, err)
		}

		if c.ExpectedTime != actualTime {
			t.Errorf("Case: %v, expected report time to be %v, but was %v", l, c.ExpectedTime, actualTime)
		}

		if err == nil && c.ExpectedErr {
			t.Errorf("Case: %v, missing expected error", l)
		}
		if err != nil && !c.ExpectedErr {
			t.Errorf("Case: %v, unexpected error: %v", l, err)
		}

	}
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(FILENAME_TIMESTAMP_LAYOUT, timeStr)
	return t
}

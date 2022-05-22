package core

import (
	"testing"
	"time"
)

// ResourcesReport tests

func TestResourcesReportItemEquivalent(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		A ResourcesReportItem
		B ResourcesReportItem
		E bool
	}{
		"Zero values":                  {ResourcesReportItem{}, ResourcesReportItem{}, true},
		"Same time":                    {ResourcesReportItem{AnalysisTime: now}, ResourcesReportItem{AnalysisTime: now}, true},
		"Different Times":              {ResourcesReportItem{AnalysisTime: now}, ResourcesReportItem{AnalysisTime: now.Add(time.Minute)}, true},
		"Different Name":               {ResourcesReportItem{ResourceName: `a`}, ResourcesReportItem{ResourceName: `b`}, false},
		"Different ARN":                {ResourcesReportItem{ResourceARN: `a`}, ResourcesReportItem{ResourceARN: `b`}, false},
		"Different Type":               {ResourcesReportItem{ResourceType: `a`}, ResourcesReportItem{ResourceType: `b`}, false},
		"Different TagBU":              {ResourcesReportItem{ResourceTagBusinessUnit: `a`}, ResourcesReportItem{ResourceTagBusinessUnit: `b`}, false},
		"Different TagEnv":             {ResourcesReportItem{ResourceTagEnvironment: `a`}, ResourcesReportItem{ResourceTagEnvironment: `b`}, false},
		"Different TagOwner":           {ResourcesReportItem{ResourceTagOwner: `a`}, ResourcesReportItem{ResourceTagOwner: `b`}, false},
		"Different TagConfidentiality": {ResourcesReportItem{ResourceTagConfidentiality: `a`}, ResourcesReportItem{ResourceTagConfidentiality: `b`}, false},
		"Different TagIntegrity":       {ResourcesReportItem{ResourceTagIntegrity: `a`}, ResourcesReportItem{ResourceTagIntegrity: `b`}, false},
		"Different TagAvailability":    {ResourcesReportItem{ResourceTagAvailability: `a`}, ResourcesReportItem{ResourceTagAvailability: `b`}, false},
		"Different Tags":               {ResourcesReportItem{ResourceTags: `a`}, ResourcesReportItem{ResourceTags: `b`}, false},
	}
	for l, c := range cases {
		if o := c.A.Equivalent(c.B); o != c.E {
			t.Errorf("Case: %v", l)
		}
	}
}

func TestDecodeResourcesReportItem(t *testing.T) {
	// 2021-06-11T20:54:08.112773+00:00,AccountAdminAccessRole-Sandbox,arn:aws:iam::139710491120:role/AccountAdminAccessRole-Sandbox,IAMRole,,,,,,,{}
	// validTime, _ := time.Parse(time.RFC3339Nano, `2021-06-11T20:54:08.112773+00:00`)
	cases := map[string]struct {
		Fields      []string
		Expected    ResourcesReportItem
		ExpectedErr bool
	}{
		`Basic working deserialization`: {
			Fields: []string{
				`2021-06-11T20:54:08.112773+00:00`,
				`AccountAdminAccessRole-Sandbox`,
				`arn:aws:iam::139710491120:role/AccountAdminAccessRole-Sandbox`,
				`IAMRole`, `a`, `b`, `c`, `d`, `e`, `f`, `{}`},
			Expected: ResourcesReportItem{
				ResourceName:               `AccountAdminAccessRole-Sandbox`,
				ResourceARN:                `arn:aws:iam::139710491120:role/AccountAdminAccessRole-Sandbox`,
				ResourceType:               `IAMRole`,
				ResourceTagBusinessUnit:    `a`,
				ResourceTagEnvironment:     `b`,
				ResourceTagOwner:           `c`,
				ResourceTagConfidentiality: `d`,
				ResourceTagIntegrity:       `e`,
				ResourceTagAvailability:    `f`,
				ResourceTags:               `{}`,
			},
			ExpectedErr: false,
		},
		`Bad date`: {
			Fields: []string{
				`fdjfkldjsakfldsjkl`,
				`AccountAdminAccessRole-Sandbox`,
				`arn:aws:iam::139710491120:role/AccountAdminAccessRole-Sandbox`,
				`IAMRole`, `a`, `b`, `c`, `d`, `e`, `f`, `{}`},
			Expected:    ResourcesReportItem{},
			ExpectedErr: true,
		},
	}
	for l, c := range cases {
		o, err := DecodeResourcesReportItem(c.Fields)
		if !c.Expected.Equivalent(o) {
			t.Errorf("Case: %v, output does not match expectation %v vs %v", l, o, c.Expected)
		}
		if err == nil && c.ExpectedErr {
			t.Errorf("Case: %v, missing expected error", l)
		}
		if err != nil && !c.ExpectedErr {
			t.Errorf("Case: %v, unexpected error: %v", l, err)
		}
	}
}

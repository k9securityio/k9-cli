package cmd

// version should be set to the latest git tag using ldflags
var version string
var revision string
var buildtime string

const (
	// EnvPrefix defines the prefix that this program uses to distinguish its environment variables.
	EnvPrefix = `K9`
)

const (
	FLAG_CUSTOMER_ID   = `customer_id`
	FLAG_VERBOSE       = `verbose`
	FLAG_ACCOUNT       = `account`
	FLAG_FORMAT        = `format`
	FLAG_ANALYSIS_DATE = `analysis-date`
	FLAG_REPORT_HOME   = `report-home`

	FLAG_ARN  = `arn`
	FLAG_ARNS = `arns`

	FLAG_NAME  = `name`
	FLAG_NAMES = `names`
)

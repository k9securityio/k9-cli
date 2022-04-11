package cmd

// version should be set to the latest git tag using ldflags
var version string
var revision string
var buildtime string

const (
	// EnvPrefix defines the prefix that this program uses to distinguish its environment variables.
	EnvPrefix = `K9`
)

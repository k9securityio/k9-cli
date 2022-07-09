package cmd

import (
	"fmt"
	"io"

	"github.com/k9securityio/k9-cli/core"
)

func DumpDBStats(o io.Writer, db *core.DB) {
	customers, accounts, total := db.Sizes()
	fmt.Fprintf(o, "Local database:\n\tCustomers:\t\t%v\n\tAccounts:\t\t%v\n\tTotal analysis dates: \t%v\n",
		customers, accounts, total)
}

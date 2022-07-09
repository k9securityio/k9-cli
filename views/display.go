package views

import (
	"encoding/json"
	"fmt"
	"io"
)

func Display(stdout, stderr io.Writer, format string, report interface{}) {
	switch format {
	case `pdf`:
	case `csv`:
		WriteCSVTo(stdout, stderr, report)
	case `tap`:
	case `json`:
		b, err := json.Marshal(report)
		if err != nil {
			fmt.Fprintln(stderr, `unable to marshal report to json`)
		}
		fmt.Fprintln(stdout, string(b))
	default:
		fmt.Fprintln(stderr, `invalid output type`)
	}
}

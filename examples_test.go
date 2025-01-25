package timeserie_test

import (
	"os"
	"strings"

	"github.com/etnz/timeserie"
)

func Example_Load() {

	source := `
	{ "on":"01-1-2", "ts1":2, "ts2":4}
	{ "on":"01-1-1", "ts1":1, "ts2":1}
	`
	dict := make(map[string]*timeserie.Support)
	if err := timeserie.Load(dict, strings.NewReader(source)); err != nil {
		panic(err)
	}
	// and print it out
	if err := timeserie.Format(os.Stdout, dict); err != nil {
		panic(err)
	}

	//Output:
	// { "on":"01-1-1", "ts2":1, "ts1":1}
	// { "on":"01-1-2", "ts2":4, "ts1":2}
}

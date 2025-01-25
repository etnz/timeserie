package timeserie

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

const (
	timeFormat = "06-1-2"
	attrOn     = "on"
)

// jsonLoad load content from a jsonlist file format.
func jsonLoad(r io.Reader) (list []any, src []string, err error) {
	line := 0
	scanner := bufio.NewScanner(r)
	list = make([]any, 0, 100000)
	for scanner.Scan() {
		line++
		var a any

		data := scanner.Bytes()

		// ignore json parsing of empty lines
		if len(strings.Trim(string(data), " \t\r")) == 0 {
			list = append(list, a)
			src = append(src, "")
			continue
		}
		src = append(src, string(data))

		err = json.Unmarshal(data, &a)
		if err != nil {
			return nil, src, fmt.Errorf("json line parsing error: line %v %q: %w", line, scanner.Text(), err)
		}
		list = append(list, a)
	}
	return list, src, nil
}

// Open supports from a value change dump file.
func Open(dict map[string]*Support, filenames ...string) (map[string]*Support, error) {
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return dict, fmt.Errorf("cannot open file %q: %w", filename, err)
		}
		defer f.Close()
		err = Load(dict, f)
		if err != nil {
			return nil, fmt.Errorf("cannot read %q content: %w", filename, err)
		}
	}
	return dict, nil
}

// Load support from a value change dump stream.
func Load(dict map[string]*Support, r io.Reader) error {
	// Read the source.
	lines, src, err := jsonLoad(r)
	if err != nil {
		return fmt.Errorf("load support error: %w", err)
	}

	for i, line := range lines {
		jline := src[i] // For error creation.
		if line == nil {
			// simply ignore empty lines.
			continue
		}
		jmap, ok := line.(map[string]any)
		if !ok {
			return fmt.Errorf("load support error line %v: json object is required but got %q", i, jline)
		}
		jstring, ok := jmap[attrOn]
		if !ok {
			return fmt.Errorf("load support error line %v: json object is missing the attribute 'on' with a date: %q", i, jline)
		}
		if _, ok := jstring.(string); !ok {
			return fmt.Errorf("load support error line %v: attribute 'on' must be of type 'string' got %q", i, jline)
		}

		on, err := time.Parse(timeFormat, jstring.(string))
		if err != nil {
			return fmt.Errorf("load support error line %v: attribute 'on' must be a valid date in the format %q got %q", i, timeFormat, jstring)
		}
		// Read all other attributes as (key,value) pairs of (Timeserie name, Timeserie value)
		for id, quantity := range jmap {
			if id == attrOn { // skip this one
				// reserved word for timestamp
				continue
			}
			_, ok := quantity.(float64)
			if !ok {
				return fmt.Errorf("load support error line %v: attribute %q must be a valid number got %q", i, id, jline)
			}
			s, ok := dict[id]
			if !ok {
				s = new(Support)
				dict[id] = s
			}
			s.Append(on, quantity.(float64))
		}
	}
	return nil
}

func Format(w io.Writer, dict map[string]*Support) error {
	var fs []*Function
	var ids []string
	for k, v := range dict {
		fs = append(fs, New(v, ModeNullset))
		ids = append(ids, k)
	}
	for t := range Iterate(fs...) {
		_, err := fmt.Fprintf(w, "{ %q:%q", attrOn, t.Format(timeFormat))
		if err != nil {
			return err
		}
		for i, id := range ids {
			f := fs[i]
			v := f.F(t)
			if !math.IsNaN(v) {
				_, err := fmt.Fprintf(w, ", %q:%v", id, v)
				if err != nil {
					return err
				}
			}
		}
		_, err = fmt.Fprintln(w, "}")
		if err != nil {
			return err
		}
	}
	return nil
}

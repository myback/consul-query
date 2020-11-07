/*
Copyright © 2020 MyBack

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"os"
	"text/tabwriter"
)

type outputFormat uint8
type outputData map[string][]string

const (
	textFormat outputFormat = iota
	jsonFormat
	jsonPrettyFormat
)

func (i *outputFormat) String() string {
	switch *i {
	case jsonFormat:
		return "json"
	case jsonPrettyFormat:
		return "jsonp"
	default:
		return "text"
	}
}

func (i *outputFormat) Set(s string) error {
	switch s {
	case "json":
		*i = jsonFormat
	case "jsonp":
		*i = jsonPrettyFormat
	case "text":
		*i = textFormat
	default:
		return fmt.Errorf("unknown output format: %s", s)
	}

	return nil
}

func (i *outputFormat) Type() string {
	return "outputFormat"
}

func anyTag(a, b []string) bool {
	for _, _a := range a {
		for _, _b := range b {
			if _a == _b {
				return true
			}
		}
	}

	return false
}

func allTags(a, b []string) bool {
	set := map[string]bool{}
	for _, v := range a {
		set[v] = true
	}

	for _, v := range b {
		if _, ok := set[v]; !ok {
			return false
		}
	}

	return true
}

func errCheck(e error, code int) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(code)
	}
}

func consulClient() (*consulApi.Client, error) {
	return consulApi.NewClient(consulApi.DefaultConfig())
}

func fprint(minwidth int, o interface{}, d outputData) {
	switch output {
	case jsonFormat:
		jsonPrint(o, false)
	case jsonPrettyFormat:
		jsonPrint(o, true)
	default:
		tabwidth := 8
		if minwidth >= tabwidth {
			minwidth -= tabwidth
		}
		tw := tabwriter.NewWriter(os.Stdout, minwidth, tabwidth, 1, '\t', tabwriter.AlignRight)

		for head, data := range d {
			if !noHeaders {
				fmt.Fprintln(tw, head)
			}
			for _, d := range data {
				fmt.Fprintln(tw, d)
			}
		}

		tw.Flush()
	}
}

func jsonPrint(i interface{}, pretty bool) {
	enc := json.NewEncoder(os.Stdout)
	if pretty {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(i); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

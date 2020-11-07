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
	"github.com/spf13/cobra"
)

var (
	output    outputFormat
	noHeaders bool

	rootCmd = &cobra.Command{
		Use:   "consul-query",
		Short: "Simple query to Hashicorp Consul",
		Long: `Simple query to Hashicorp Consul`,
	}
)

// Execute executes the root command.
func Execute() {
	errCheck(rootCmd.Execute(), 1)
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&noHeaders, "no-headers", "n", false, "don't show headers. Text format only")
	rootCmd.PersistentFlags().VarP(&output, "output", "o", "output format: json, jsonp, text")

	rootCmd.AddCommand(keyValueCmd)
	rootCmd.AddCommand(serviceCmd)
}

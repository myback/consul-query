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
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"strings"
)

var (
	tags []string
	anyT bool
)

var serviceCmd = &cobra.Command{
	Use:   "svc [service name]",
	Short: "List Consul services",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			errCheck(getServices(tags), 3)
			return
		}

		errCheck(getCatalog(args[0], tags), 5)
	},
}

func init() {
	serviceCmd.PersistentFlags().StringArrayVarP(&tags, "tag", "t", nil, "list of tags")
	serviceCmd.PersistentFlags().BoolVar(&anyT, "any", false, "any tag match")
}

func getServices(tagsList []string) error {
	client, err := consulClient()
	if err != nil {
		return err
	}

	lst, _, err := client.Catalog().Services(&consulApi.QueryOptions{})
	if err != nil {
		return err
	}

	const head = "SERVICE\tTAGS"
	data := outputData{}
	compareTags := allTags
	if anyT {
		compareTags = anyTag
	}

	w := 0
	for k, svcTags := range lst {
		svcTagsOut := svcTags
		if len(svcTagsOut) > 10 {
			svcTagsOut = svcTagsOut[:10]
		}

		if len(tagsList) > 0 {
			if compareTags(svcTags, tagsList) {
				if len(k) > w {
					w = len(k)
				}
				data[head] = append(data[head], fmt.Sprintf("%s\t%s", k, strings.Join(svcTagsOut, ":")))
			}
		} else {
			if len(k) > w {
				w = len(k)
			}
			data[head] = append(data[head], fmt.Sprintf("%s\t%s", k, strings.Join(svcTagsOut, ":")))
		}
	}

	fprint(w, lst, data)

	return nil
}

func getCatalog(svcName string, tagsList []string) error {
	client, err := consulClient()
	if err != nil {
		return err
	}

	lst, _, err := client.Catalog().ServiceMultipleTags(svcName, tagsList, &consulApi.QueryOptions{})
	if err != nil {
		return err
	}

	const head = "NODE\tADDRESS\tTAGS"
	data := outputData{}
	w := 0
	for _, svc := range lst {
		if len(svc.Node) > w {
			w = len(svc.Node)
		}
		data[head] = append(data[head], fmt.Sprintf("%s\t%s:%d\t%s", svc.Node, svc.Address, svc.ServicePort, strings.Join(svc.ServiceTags, ":")))
	}

	fprint(w, lst, data)

	return nil
}

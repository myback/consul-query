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
	recursive bool
)

var keyValueCmd = &cobra.Command{
	Use:   "kv key",
	Short: "Query to Consul key-value",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			errCheck(listKeys(""), 3)
			return
		}

		if strings.HasSuffix(args[0], "/") && !recursive {
			errCheck(listKeys(args[0]), 5)
		} else {
			errCheck(getKV(args[0], recursive), 5)
		}
	},
}

func init() {
	keyValueCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "Recursive get keys")
}

func listKeys(pfx string) error {
	client, err := consulClient()
	if err != nil {
		return err
	}

	pfx = strings.TrimPrefix(pfx, "/")
	if !strings.HasSuffix(pfx, "/") && pfx != "" {
		pfx += "/"
	}

	lst, _, err := client.KV().Keys(pfx, "", &consulApi.QueryOptions{})
	if err != nil {
		return err
	}

	const head = "KEYS"
	data := outputData{}

	keys := map[string]bool{}
	w := 0
	for _, k := range lst {
		key := strings.Split(strings.TrimPrefix(k, pfx), "/")
		keyGroup := pfx + key[0]
		if len(key) > 1 {
			keyGroup += "/"
		}

		if len(keyGroup) > w {
			w = len(keyGroup)
		}

		if ok := keys[keyGroup]; !ok {
			keys[keyGroup] = true
			data[head] = append(data[head], keyGroup)
		}
	}

	fprint(w, lst, data)

	return nil
}

func getKV(key string, recurse bool) error {
	client, err := consulClient()
	if err != nil {
		return err
	}

	const head = "KEY\tVALUE"
	data := outputData{}

	if recurse {
		kvList, _, err := client.KV().List(key, &consulApi.QueryOptions{})
		if err != nil {
			return err
		}

		w := 0
		for _, kv := range kvList {
			if len(kv.Key) > w {
				w = len(kv.Key)
			}
			data[head] = append(data[head], fmt.Sprintf("%s\t%s", kv.Key, kv.Value))
		}

		fprint(w, kvList, data)

		return nil
	}

	kvPair, _, err := client.KV().Get(key, &consulApi.QueryOptions{})
	if err != nil {
		return err
	}

	if kvPair == nil {
		return nil
	}

	data[head] = append(data[head], fmt.Sprintf("%s\t%s", kvPair.Key, kvPair.Value))

	fprint(0, kvPair, data)

	return nil
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	consulApi "github.com/hashicorp/consul/api"
)

type outputFormat uint8
type argList []string

const (
	textFormat outputFormat = iota
	jsonFormat
	jsonPrettyFormat
)

func (i *argList) String() string {
	if len(*i) == 0 {
		return ""
	}

	return fmt.Sprint(*i)
}

func (i *argList) Set(s string) error {
	*i = append(*i, s)
	return nil
}

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

func main() {
	var (
		noHeaders   bool
		output      outputFormat
		serviceName string
		tagList     argList
	)

	flag.BoolVar(&noHeaders, "n", false, "don't show headers")
	flag.Var(&output, "o", "``output format: json, jsonp, text (default text)")
	flag.Var(&tagList, "t", "``list of tags")

	flag.Usage = flagUsage

	// run program without arguments
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	// first arg is consul service
	serviceName = os.Args[1]

	if strings.HasPrefix(serviceName, "-") {
		// try find -h argument
		flag.Parse()

		// if not found -h
		flag.Usage()
		os.Exit(1)
	}

	flag.CommandLine.Parse(os.Args[2:])

	ls, err := consulQuery(serviceName, tagList)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fprint(os.Stdout, ls, output, noHeaders)
}

func flagUsage() {
	fmt.Printf("Usage: %s service [arguments ...]\n  Arguments:\n", os.Args[0])
	flag.PrintDefaults()
}

func consulQuery(svc string, tags []string) ([]*consulApi.CatalogService, error) {
	client, err := consulApi.NewClient(consulApi.DefaultConfig())
	if err != nil {
		return nil, err
	}

	svcList, _, err := client.Catalog().ServiceMultipleTags(svc, tags, &consulApi.QueryOptions{})
	if err != nil {
		return nil, err
	}

	return svcList, nil
}

func fprint(w io.Writer, svcList []*consulApi.CatalogService, f outputFormat, noHead bool) {
	switch f {
	case jsonFormat:
		if err := json.NewEncoder(w).Encode(svcList); err != nil {
			fmt.Println(err)
			os.Exit(3)
		}

	case jsonPrettyFormat:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(svcList); err != nil {
			fmt.Println(err)
			os.Exit(3)
		}

	default:
		tw := tabwriter.NewWriter(w, 0, 8, 1, '\t', tabwriter.AlignRight)
		if !noHead {
			fmt.Fprintln(tw, "NODE\tADDRESS\tTAGS")
		}

		for _, svc := range svcList {
			fmt.Fprintf(tw, "%s\t%s:%d\t%s\n", svc.Node, svc.Address, svc.ServicePort, strings.Join(svc.ServiceTags, ":"))
		}
		tw.Flush()
	}
}

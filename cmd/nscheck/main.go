package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/dogasantos/nscheck/pkg/runner/nscheck"
)

type Options struct {
	ResolverFile		string
	TrustedNs			string
	Version				bool
	Verbose				bool
}

var version = "0.1"

func parseOptions() *Options {
	options := &Options{}
	flag.StringVar(&options.ResolverFile, 		"l", "", "List of dns servers we should test")
	flag.StringVar(&options.TrustedNs, 			"s", "", "Set a specific IP as the trusted server, otherwise, use default (Google and Cloudflare)")
	flag.BoolVar(&options.Version, 				"i", false, "Version info")
	flag.BoolVar(&options.Verbose, 				"v", false, "Verbose mode")
	flag.Parse()
	return options
}

func main() {
	var wg sync.WaitGroup
	options := parseOptions()
	if options.Version {
		fmt.Println(version)
	}
	
	if options.ResolverFile != "" {
		if options.Verbose == true {
			fmt.Printf("[+] NSCHECK v%s\n",version)
		}
		ResolverFilestream, _ := ioutil.ReadFile(options.ResolverFile)
		resolversContent := string(ResolverFilestream)
		resolvers := strings.Split(resolversContent, "\n")

		if options.Verbose == true {
			fmt.Printf("[*] Resolvers loaded: %d\n",len(resolvers))
			if len(options.TrustedNs)>0 {
				fmt.Printf("[*] Trusted NS record: %s\n",options.TrustedNs)
			}
		}
		
		for _, prefix := range resolvers {
			wg.Add(1)
			go nscheck.CheckResolvers(resolvers, options.TrustedNs, &wg, options.Verbose)
		}
		wg.Wait()
		if options.Verbose == true {
			fmt.Printf("[.] Done.\n")
		}
	} 
}





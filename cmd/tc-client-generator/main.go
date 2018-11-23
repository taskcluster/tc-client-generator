package main

import (
	"log"

	docopt "github.com/docopt/docopt-go"
)

var (
	version  = "0.0.0"
	revision = "" // this is set during build with `-ldflags "-X main.revision=$(git rev-parse HEAD)"`
)

func usage(versionName string) string {
	return versionName + `

tc-client-generator generates taskcluster clients in a variety of programming languages.

  Usage:
    tc-client-generator --help
    tc-client-generator --version

  Exit Codes:
    0      Completed successfully.
`
}

// Entry point into the generic worker...
func main() {
	versionName := "tc-client-generator " + version
	if revision != "" {
		versionName += " [ revision: https://github.com/taskcluster/tc-client-generator/commits/" + revision + " ]"
	}
	_, err := docopt.Parse(usage(versionName), nil, true, versionName, false, true)
	if err != nil {
		log.Println("Error parsing command line arguments!")
		panic(err)
	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	version  = "0.0.1"
	revision = "" // this is set during build with `-ldflags "-X main.revision=$(git rev-parse HEAD)"`
)

func usage(versionName string) string {
	return versionName + `

tc-client-generator generates taskcluster clients in a variety of programming
languages.

Usage:
  tc-client-generator [--python-dir PYTHON_DIR]
                      [--node-dir NODE_DIR]
                      [--go-dir GO_DIR]
                      [--javascript-dir JAVASCRIPT_DIR]
                      [--java-dir JAVA_DIR]
                      [--taskcluster-root-url TASKCLUSTER_ROOT_URL]
  tc-client-generator -h|--help
  tc-client-generator --version

If called without any arguments, tc-client-generator will generate _all_
clients in subdirectories of the current directory. Otherwise, only language
clients that have a directory specified as a command option will be generated.

If --taskcluster-root-url is not specified, the environment variable
TASKCLUTER_ROOT_URL will be used if set, otherwise no changes will be made and
exit code 64 will be returned.

Examples:

1) Generate all language clients in subdirectories of the current directory:

  $ tc-client-generator

2) Generate go client in ~/go/src/github.com/foo/bar/vendor and python client
in current directory:

  $ tc-client-generator --go-dir ~/go/src/github.com/foo/bar/vendor --python-dir .

Exit Codes:

  0     Completed successfully.
  1     Invalid arguments passed.
 64     No taskcluster root URL specified.
 65     Invalid taskcluster root url specified.
 66     Error fetching references/manifests from taskcluster.
 67     Error creating python source directory.
 68     Error creating node source directory.
 99     Error creating go source directory.
 70     Error creating javascript source directory.
 71     Error creating java source directory.
 72     Error writing to the filesystem.
 73     Internal error (crash).
`
}

type GenerateClients struct {
	PythonDir     string
	NodeDir       string
	GoDir         string
	JavascriptDir string
	JavaDir       string
}

func main() {
	// Pass in command line args and stderr/stdout explicitly to
	// ProcessArguments so that tests can pass in test/mock objects.
	exitCode := ProcessArguments(os.Args, os.Stdout, log.New(os.Stderr, "tc-client-generator: ", 0))
	os.Exit(exitCode)
}

// ProcessArguments takes a set of command line arguments to process, an
// io.Writer to write its output to, and a logger for writing error messages
// to. It performs all actions required for the arguments passed, and returns
// an exit code for the process. The arguments are passed in, rather than being
// implicitly defined within the function, in order that tests can pass in test
// command line arguments, a mock output io.Writer and logger, and evaluate the
// outputs and exit code of the function.
//
// Note, I wasn't able to have docopt successfully reject additional unwanted
// arguments, and therefore I have implemented argument parsing directly.
func ProcessArguments(args []string, out io.Writer, logger *log.Logger) (exitCode int) {

	defer func() {
		if exitCode != 0 {
			logger.Printf("Exiting with code %v", exitCode)
		}
	}()

	// Assuming the binary was built correctly, the linker should have included the git revision when linking,
	// so build the version name at runtime. Gracefully omit the revision if the build doesn't contain one.
	versionName := "tc-client-generator " + version
	if revision != "" {
		versionName += " [ revision: https://github.com/taskcluster/tc-client-generator/commits/" + revision + " ]"
	}

	// Check for standard options that don't generate a client (-h|--help|--version).
	switch {
	case len(args) == 0:
		panic("Somehow the command arguments have been lost - this should not be possible and indicates a bug.")
	case len(args) == 2 && (args[1] == "--help" || args[1] == "-h"):
		fmt.Fprintln(out, usage(versionName))
		return
	case len(args) == 2 && args[1] == "--version":
		fmt.Fprintln(out, versionName)
		return
	}

	// From this point on, we assume that the args are to create one or more clients.

	// Clients to generate
	var clients GenerateClients
	var taskclusterRootURL string
	noLanguagesSpecified := true

	// Consume the arguments by looping through them. If a parameter consumes
	// more than one argument, `i` is incremented inside its handler, in order
	// that at the start of the loop, args[i] should always be a recognised
	// command parameter.
	for i := 1; i < len(args); i++ {
		var err error
		switch args[i] {
		case "--python-dir":
			noLanguagesSpecified = false
			err = parseClientOptions(args, &i, &clients.PythonDir)
		case "--node-dir":
			noLanguagesSpecified = false
			err = parseClientOptions(args, &i, &clients.NodeDir)
		case "--go-dir":
			noLanguagesSpecified = false
			err = parseClientOptions(args, &i, &clients.GoDir)
		case "--javascript-dir":
			noLanguagesSpecified = false
			err = parseClientOptions(args, &i, &clients.JavascriptDir)
		case "--java-dir":
			noLanguagesSpecified = false
			err = parseClientOptions(args, &i, &clients.JavaDir)
		case "--taskcluster-root-url":
			if i+1 >= len(args) {
				logger.Printf("No taskcluster root url specified after `--taskcluster-root-url` argument. This option requires a further argument.")
				exitCode = 1
				return
			}
			i++
			taskclusterRootURL = args[i]
		default:
			logger.Printf("Unrecognised option: %v", args[i])
			exitCode = 1
			return
		}
		if err != nil {
			logger.Printf("%v", err)
			exitCode = 1
			return
		}
	}

	// fall back to using environment variable TASKCLUSTER_ROOT_URL if no
	// --taskcluster-root-url specified
	if taskclusterRootURL == "" {
		taskclusterRootURL = os.Getenv("TASKCLUSTER_ROOT_URL")
		if taskclusterRootURL == "" {
			logger.Print("No taskcluster root URL specified. Please provide option `--taskcluster-root-url` or set environment variable TASKCLUSTER_ROOT_URL.")
			exitCode = 64
			return
		}
	}

	// If no languages were specified, then all clients should be generated in
	// default locations.
	if noLanguagesSpecified {
		clients.PythonDir = "generated-clients/python"
		clients.NodeDir = "generated-clients/node"
		clients.GoDir = "generated-clients/go"
		clients.JavascriptDir = "generated-clients/javascript"
		clients.JavaDir = "generated-clients/java"
	}

	fmt.Fprintf(out, "Building the following clients: %#v\n", clients)
	fmt.Fprintf(out, "Using taskcluster root URL %v\n", taskclusterRootURL)
	return
}

// parseClientOptions attempts to
func parseClientOptions(args []string, index *int, directory *string) error {
	if *index+1 >= len(args) {
		return fmt.Errorf("No directory specified after parameter %v - this parameter requires an argument.", args[*index])
	}
	*index++
	if *directory != "" {
		return fmt.Errorf("Two directories specified for %v: %q and %q - maximum of one allowed.", args[*index-1], *directory, args[*index])
	}
	*directory = args[*index]
	return nil
}

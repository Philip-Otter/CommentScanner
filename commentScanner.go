/*
An application to quickly pull out comments and other information from a webpage
2xdropout 2024
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

type target struct {
	URL           string
	FindEmail     bool
	FindCreds     bool
	FindPhone     bool
	FindSources   bool
	FindRefs      bool
	MaxDepth      int
	UserAgent     string
	ReportingMode string

	IsOutput     bool `default:"false"`
	CurrentDepth int  `default:"1"`
}

func search(targetptr *target, workerptr *int, maxWorkers int) {
	fmt.Println("Searching")
}

func main() {
	title := `
	╔═╗┌─┐┌┬┐┌┬┐┌─┐┌┐┌┌┬┐  ╔═╗┌─┐┌─┐┌┐┌┌┐┌┌─┐┬─┐
	║  │ │││││││├┤ │││ │   ╚═╗│  ├─┤││││││├┤ ├┬┘
	╚═╝└─┘┴ ┴┴ ┴└─┘┘└┘ ┴   ╚═╝└─┘┴ ┴┘└┘┘└┘└─┘┴└─
				by The 2xdropout
					2024`
	seperator := "============================================================="

	// required flags
	urlFlagptr := flag.String("u", "", "The target URL")
	// Optional content searching flags
	emailFlagptr := flag.Bool("e", false, "Search for emails")
	credentialFlagptr := flag.Bool("c", false, "Search for credentials")
	phoneFlagptr := flag.Bool("p", false, "Search for credentials")
	sourceFlageptr := flag.Bool("s", false, "Search for source files")
	referenceFlagptr := flag.Bool("r", false, "Search for references")
	// Configuration flags
	depthFlagptr := flag.Int("depth", 1, "Link and reference scanning depth")
	workersFlagptr := flag.Int("workers", 10, "Concurrent scanning and parsing workers")
	userAgentFlagptr := flag.String("userAgent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0", "Set the user agent with which requests are made")
	reportingModeFlagptr := flag.String("reporting", "stdout", "stdout, file, tree, map, html, fuzzFiles")

	flag.Parse()

	// Print running config
	color.Green(title + "\n" + seperator)
	color.Green("URL:  			" + *urlFlagptr)
	if *emailFlagptr {
		color.Green("Find Emails:  		" + strconv.FormatBool(*emailFlagptr))
	}
	if *credentialFlagptr {
		color.Green("Find Credentials:  	" + strconv.FormatBool(*credentialFlagptr))
	}
	if *phoneFlagptr {
		color.Green("Find Phone Numbers:  	" + strconv.FormatBool(*phoneFlagptr))
	}
	if *sourceFlageptr {
		color.Green("Find Source Files:  	" + strconv.FormatBool(*sourceFlageptr))
	}
	if *referenceFlagptr {
		color.Green("Find Reference Files:  	" + strconv.FormatBool(*referenceFlagptr))
	}

	color.Green("Number Of Workers:  	" + strconv.Itoa(*workersFlagptr))
	color.Green("Scanning Depth:  	" + strconv.Itoa(*depthFlagptr))
	color.Green("Reporting Mode:  	" + *reportingModeFlagptr)
	color.Green(seperator)

	if *urlFlagptr == "" {
		color.Red("URL flag not set!")
		os.Exit(2)
	}

	currentWorkers := 1
	currentWorkersptr := &currentWorkers

	rootTarget := target{
		URL:           *urlFlagptr,
		FindEmail:     *emailFlagptr,
		FindCreds:     *credentialFlagptr,
		FindPhone:     *phoneFlagptr,
		FindSources:   *sourceFlageptr,
		FindRefs:      *referenceFlagptr,
		MaxDepth:      *depthFlagptr,
		UserAgent:     *userAgentFlagptr,
		ReportingMode: *reportingModeFlagptr,
	}

	fmt.Println(rootTarget)

	search(&rootTarget, currentWorkersptr, *workersFlagptr)
}

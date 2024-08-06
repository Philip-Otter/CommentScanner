/*
An application to quickly pull out comments and other information from a webpage
2xdropout 2024
*/
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
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
	DisableBasic  bool
	MaxDepth      int
	UserAgent     string
	ReportingMode string

	IsOutput     bool `default:"false"`
	CurrentDepth int  `default:"1"`
}

func get_HTML_Comments(textBlock string) {
	commentRegex, _ := regexp.Compile(`(?s)\<\!\-\-.+?\-\-\>`)
	commentList := commentRegex.FindAllString(textBlock, -1)

	fmt.Print("| ")
	color.Yellow("HTML COMMENTS:  ")
	for _, item := range commentList {
		fmt.Print("|#")
		color.Cyan(item)
	}
}

func get_CSS_Comments(textBlock string) {
	commentRegex, _ := regexp.Compile(`(?s)\/\*.+?\*/`)
	commentList := commentRegex.FindAllString(textBlock, -1)

	fmt.Print("| ")
	color.Yellow("CSS COMMENTS:  ")
	for _, item := range commentList {
		fmt.Print("|#")
		color.Cyan(item)
	}
}

func get_JS_Comments(textBlock string) {
	inlineCommentRegex, _ := regexp.Compile(`[^:\\]\/\/.+`)
	inlineCommentList := inlineCommentRegex.FindAllString(textBlock, -1)

	blockCommentRegex, _ := regexp.Compile(`(?s)\/\*.+?\*\/`)
	blockCommentList := blockCommentRegex.FindAllString(textBlock, -1)

	fmt.Print("| ")
	color.Yellow("JS COMMENTS:  ")
	for _, item := range inlineCommentList {
		fmt.Print("|#")
		color.Cyan(item)
	}
	for _, item := range blockCommentList {
		fmt.Print("|#")
		color.Cyan(item)
	}
}

func get_emails(textBlock string) {
	emailRegex, _ := regexp.Compile(`[a-zA-Z0-9#$%&'*+-/=?^_\|{}~]+?\@[a-zA-Z0-9\.]+\.[a-z]+[^\s0-9"<!@#$%^&*(){}>;',.]`)
	emailList := emailRegex.FindAllString(textBlock, -1)

	fmt.Print("| ")
	color.Yellow("Possible Emails")
	for _, item := range emailList {
		fmt.Print("|#")
		color.Cyan(item)
	}
}

func search(targetptr *target, workerptr *int, maxWorkers int) {
	seperatorString := "----------------------------------------------------------------------------------------"
	fmt.Println(seperatorString)
	color.Red(*&targetptr.URL)
	fmt.Println(seperatorString)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Get(*&targetptr.URL)
	if err != nil {
		log.Print(err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
	}

	rawFile := string(body)

	if !*&targetptr.DisableBasic {
		get_HTML_Comments(rawFile)
		get_CSS_Comments(rawFile)
		get_JS_Comments(rawFile)
	}
	if *&targetptr.FindEmail {
		get_emails(rawFile)
	}
}

func main() {
	title := `
	╔═╗┌─┐┌┬┐┌┬┐┌─┐┌┐┌┌┬┐  ╔═╗┌─┐┌─┐┌┐┌┌┐┌┌─┐┬─┐
	║  │ │││││││├┤ │││ │   ╚═╗│  ├─┤││││││├┤ ├┬┘
	╚═╝└─┘┴ ┴┴ ┴└─┘┘└┘ ┴   ╚═╝└─┘┴ ┴┘└┘┘└┘└─┘┴└─
				by The 2xdropout
					2024`
	seperator := "========================================================================================"

	// required flags
	urlFlagptr := flag.String("u", "", "The target URL")
	// Optional content searching flags
	emailFlagptr := flag.Bool("e", false, "Search for emails")
	credentialFlagptr := flag.Bool("c", false, "Search for credentials")
	phoneFlagptr := flag.Bool("p", false, "Search for credentials")
	sourceFlageptr := flag.Bool("s", false, "Search for source files")
	referenceFlagptr := flag.Bool("r", false, "Search for references")
	noBasicFlagptr := flag.Bool("noBasic", false, "Disable basic comment scanning (HTML/CSS/JS)")
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
	if *noBasicFlagptr {
		color.Green("No HTML/CSS/JS Scan:	" + strconv.FormatBool(*noBasicFlagptr))
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
		DisableBasic:  *noBasicFlagptr,
		MaxDepth:      *depthFlagptr,
		UserAgent:     *userAgentFlagptr,
		ReportingMode: *reportingModeFlagptr,
	}

	search(&rootTarget, currentWorkersptr, *workersFlagptr)
}

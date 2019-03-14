package main

import (
	"flag"
	"fmt"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"os"
	"strings"
	"time"
)

var configPath string
var outputPath string
var format string
var verbose bool

var databasePath = "/var/lib/autopeer"

func init(){
	defer flushTelemetry()
	var err error

	// load config file
	flag.StringVar(&configPath, "config", "-", "config file")
	flag.StringVar(&outputPath, "output", "-", "output folder")
	flag.StringVar(&format, "format", "", "override output format")
	flag.BoolVar(&verbose, "verbose", false, "no shutup")
	flag.Parse()

	// load Application Insights
	client = appinsights.NewTelemetryClient("0bb67684-cece-4a13-b1ad-1be13f66b7e1")
	if verbose {
		appinsights.NewDiagnosticsMessageListener(func(msg string) error {
			fmt.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
			return nil
		})
	}

	// create database
	err = os.MkdirAll(databasePath, os.ModePerm)
	hardFail(err)
}

func main() {
	// var err error

	//pwd, err := os.Getwd()
	//hardFail(err)

	var retCode int

	action := strings.ToLower(flag.Args()[0])
	if funcPtr, ok := fnTable[action]; ok {
		retCode = funcPtr(flag.Args())
	}

	flushTelemetry()
	os.Exit(retCode)
}

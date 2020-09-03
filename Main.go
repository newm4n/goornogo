package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	coverFile       string
	minimumCoverage float64
	lineRe          = regexp.MustCompile(`^(.+):([0-9]+).([0-9]+),([0-9]+).([0-9]+) ([0-9]+) ([0-9]+)$`)
	totalStatement  = 0
	totalCounted    = 0
	lineCount       = 0
	verbose         = false
)

func main() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.StringVar(&coverFile, "i", "", "The `coverage output file` generated from - go test ./... -covermode=count -coverprofile=coverage.out")
	flag.Float64Var(&minimumCoverage, "c", -1, "The `minimum percentage` to met, value between 1 to 100")
	flag.BoolVar(&verbose, "v", true, "`verbosity` to print out coverage result, or simply emit exit code")
	flag.Parse()

	if minimumCoverage < 1 || len(coverFile) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, "Usage : -i <coverage.out> -c <minimum percentage>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	info, err := os.Stat(coverFile)
	if os.IsNotExist(err) || info.IsDir() {
		fmt.Printf("File not exist or its a directory : %s", coverFile)
		os.Exit(1)
	}

	file, err := os.Open(coverFile)
	if err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") && line != "mode: count" {
			fmt.Printf("bad mode : %s. run go test with count mode.\n", line)
			os.Exit(1)
		}
		process(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	coverage := (float64(totalCounted) / float64(totalStatement)) * float64(100)
	if coverage < minimumCoverage {
		if verbose {
			fmt.Printf("Coverage %.2f of minimum %.2f percent. [FAIL] \n", coverage, minimumCoverage)
		}
		os.Exit(1)
	} else {
		if verbose {
			fmt.Printf("Coverage %.2f of minimum %.2f percent. [OK] \n", coverage, minimumCoverage)
		}
		os.Exit(0)
	}
}

func process(line string) {
	subMatch := lineRe.FindStringSubmatch(line)
	if len(subMatch) > 0 && len(subMatch[0]) > 0 {
		lineCount++
		statements, err := strconv.Atoi(subMatch[6])
		if err != nil {
			return
		}
		count, err := strconv.Atoi(subMatch[7])
		if err != nil {
			return
		}
		totalStatement += statements
		if count > 0 {
			totalCounted += statements
		}
	}
}

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
	coverFile string
	minimumCoverage float64
	regex *regexp.Regexp
	totalStatement = 0
	totalCounted = 0
	lineCount = 0
	verbose = false
)

func main() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.StringVar(&coverFile, "i", "", "The `coverage output file` generated from - go test ./... -covermode=count -coverprofile=coverage.out")
	flag.Float64Var(&minimumCoverage, "c", -1, "The `minimum percentage` to met, value between 1 to 100")
	flag.BoolVar(&verbose, "v", true, "`verbosity` to print out coverage result, or simply emit exit code")
	flag.Parse()

	if minimumCoverage < 1 || len(coverFile) == 0 {
		_,_ = fmt.Fprintln(os.Stdout, "Usage : -i <coverage.out> -c <minimum percentage>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	regx, err := regexp.Compile(`^[a-zA-Z0-9/\.\-]+:([0-9]+).([0-9]+),([0-9]+).([0-9]+)\s+([0-9]+)\s+([0-9]+)$`)
	if err!=nil {
		panic(err)
	}
	regex=regx

	info, err := os.Stat(coverFile)
	if os.IsNotExist(err) {
		fmt.Printf("File not exist : %s", coverFile)
		os.Exit(1)
	}
	if info.IsDir() {
		fmt.Printf("Not a file : %s", coverFile)
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
	if verbose {
		fmt.Printf("Coverage %.2f of minimum %.2f percent \n", coverage, minimumCoverage)
	}
	if coverage < minimumCoverage {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func process(line string) {
	subMatch := regex.FindAllStringSubmatch(line, -1)
	if len(subMatch) > 0 && len(subMatch[0]) > 0 {
		lineCount++
		statements, err := strconv.Atoi(subMatch[0][5])
		if err != nil {
			return
		}
		count, err := strconv.Atoi(subMatch[0][6])
		if err != nil {
			return
		}
		totalStatement += statements
		if count > 0 {
			totalCounted += statements
		}
	}
}
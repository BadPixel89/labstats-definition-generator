package main

import (
	"bufio"
	"flag"
	"labstats-definition-generator/activedirectory"
	"labstats-definition-generator/excel"
	"labstats-definition-generator/jamf"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BadPixel89/colourtext"
)

var exeDir string = ""
var inFileAD = flag.String("adfile", "oulist.txt", "Specify an alternate input file containing the list of OUs. File must be: Plain text, one OU per line, located next to the executable.")
var inFileJamf = flag.String("jamffile", "jamflist.txt", "Specify an alternate input file containing the list of OUs. File must be: Plain text, one OU per line, located next to the executable.")
var outFileFlagAD = flag.String("outad", "definition-file-ad", "Specify an alternate output file for AD definitions. File will be saved next to the executable.")
var outFileFlagJamf = flag.String("outjamf", "definition-file-jamf", "Specify an alternate output file for jamf definitions. File will be saved next to the executable.")
var help = flag.Bool("h", false, "Show help text")

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir = filepath.Dir(ex)
	log.SetFlags(0) //don't log extra data such as time/date
	HandleAD()
	HandleJamf()
}

func HandleAD() {
	requestedOUs, err := ReadList(exeDir + "/" + *inFileAD)
	outFileAD := exeDir + "/" + *outFileFlagAD + ".xlsx"
	if err != nil {
		log.Fatal(err.Error())
	}

	if !ConfirmOverwriteFile(outFileAD) {
		colourtext.PrintInfo("skipping AD definition file")
		return
	}

	conn, err := activedirectory.ConnectAndBindAD()
	if err != nil {
		log.Fatal(err.Error())
	}

	results := activedirectory.SearchADbyList(conn, requestedOUs)
	os.Remove(outFileAD)
	err = excel.OutputSheetAD(outFileAD, results)
	if err != nil {
		log.Fatal(colourtext.Red + err.Error())
	}
}
func HandleJamf() {
	requestedGroups, err := ReadList(exeDir + "/" + *inFileJamf)

	outFileJamf := exeDir + "/" + *outFileFlagJamf + ".xlsx"
	if err != nil {
		log.Fatal((err.Error()))
	}
	if !ConfirmOverwriteFile((outFileJamf)) {
		colourtext.PrintInfo("skipping jamf definition file")
		return
	}
	err = jamf.JamfAuth()
	if err != nil {
		colourtext.PrintFail("Jamf auth")
	}

	results := jamf.SearchJamfByList(requestedGroups)
	os.Remove(outFileJamf)
	err = excel.OutputSheetJamf(outFileJamf, results)
	if err != nil {
		colourtext.PrintFail("writing excel file")
	}
}

// prompts user for y/n, returns true if user confirms file can be overwritten
func ConfirmOverwriteFile(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return true
	}
	reader := bufio.NewReader(os.Stdin)
	colourtext.PrintWarn("output file '" + file + "' already exists, overwrite? [y/n] : ")
	overwrite, err := reader.ReadString('\n')
	overwrite = strings.TrimSpace(overwrite)
	if err != nil {
		log.Fatal(colourtext.Red + "[exit] failed to receive console input \n" + err.Error())
	}
	switch overwrite {
	case "y":
		fallthrough
	case "Y":
		return true
	case "n":
		fallthrough
	case "N":
		colourtext.PrintInfo("use the -o flag to specify a different file name or remove/rename the existing file manually")
		return false
	default:
		colourtext.PrintInfo("please enter 'y' or 'n' to choose weather to overwrite the file. Use the -o flag to specify an outfile")
		return false
	}
}

// read in a list of OU names from the file given. 1 OU per line.
func ReadList(filepath string) ([]string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	requestedOUs := []string{}

	for scanner.Scan() {
		requestedOUs = append(requestedOUs, strings.ToUpper(scanner.Text()))
	}
	return requestedOUs, nil
}

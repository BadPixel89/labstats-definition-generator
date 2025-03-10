package excel

import (
	"fmt"
	"labstats-definition-generator/activedirectory"
	"labstats-definition-generator/jamf"
	"strings"

	"github.com/BadPixel89/colourtext"

	"github.com/xuri/excelize/v2"
)

//
//	This module will create flat hierarchy of groups based on the name of the OU as if it is one root node with all OUs beneath it.
//	If you organise your OUs in a more complex tree, LabStats can AD join and replicate this structure.
//	If you have a flat hierarchy with names that specify location, you will need to write a parsing function that generates the relevant hierarchy
//
//	The hierarchy on labstats is formed by filling out the contents of cells ascending left to right
//	Cell A1 is always the CSV of computers, A2 is the shallowest node of the tree. A3 will be the next node etc.
//

func OutputSheetAD(filepath string, ous []activedirectory.OU) error {
	File := createSheetFromOUList(ous)
	err := File.SaveAs(filepath)
	if err != nil {
		return err
	}

	colourtext.PrintDone("file generated: " + filepath)
	return nil
}

// TODO
// Set Cell B C D etc to create hierarchy - see notes at top of module and example excel file for more info
func createSheetFromOUList(ous []activedirectory.OU) *excelize.File {
	file := excelize.NewFile()
	for index, ou := range ous {
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "A", index+1), strings.Join(ou.Computers, ","))
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "B", index+1), ou.Name)
	}
	return file
}

func OutputSheetJamf(filepath string, groups []jamf.ComputerGroup) error {
	File := createSheetFromJamfList(groups)
	err := File.SaveAs(filepath)
	if err != nil {
		return err
	}
	colourtext.PrintDone("file generated: " + filepath)
	return nil
}

// TODOs
// Set Cell B C D etc to create hierarchy - see notes at top of module and example excel file for more info
func createSheetFromJamfList(groups []jamf.ComputerGroup) *excelize.File {
	file := excelize.NewFile()

	for index, group := range groups {
		compcsv := group.Computers[0].Name

		for _, computer := range group.Computers[1:] {
			compcsv += "," + computer.Name
		}
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "A", index+1), compcsv)
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "B", index+1), group.Name)
	}

	return file
}

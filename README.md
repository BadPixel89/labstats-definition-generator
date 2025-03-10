# Introduction 
Small script that takes in a list of Active Directoy Organisational Units and/or Jamf computer groups and generates a .xlsx definition file for LabStats

The benefit of this tool is that when provided with a parsing method, you can generate a tree-like hierarchy in labstats from a flat AD/Jamf structure.

If you already have a tree-like Active Directory, you could run the AD join through labstats, and it would replicate that hierarchy for you.

# Important notes

## You will need to modify this script before you use it. 

Search for "TODO" to find sections that need modification. At minimum you will need to:

* set your siteurl for jamf
* set your Active Directory LDAP address and search root

The OU and group name parsing will differ between organisations, it is intended that you create a parsing function, or slot in parsing logic, that takes in a group or OU name and creates a tree structure based on it if your naming scheme woud support it. 

See the definition file example in the repo for an idea of what labstats requires as it's input file

# Basic Use

Assuming you have cloned and appropriately modified the script:

* Create a text file containing a list of OU names called oulist.txt.
* One entry per line
* Create a text file containing a list of jamf group IDs called jamflist.txt 
* One entry per line
* place these in the same folder as the built executable
* Each line must be the name of an OU or a group ID you want to generate a definition for.
* Navigate a terminal window to the folder containing the executable and run it from the CLI (this allows you to see errors it's not strictly necessary for the tool to work)
* You should see a two new .xlsx files called have been generated in the folder next to the executable

# Build and Test
Install GoLang (This was written on version 1.22.5 - check go.mod file for this info, it should work with newer versions when they are released)

[GoLang](https://go.dev/dl/)

Clone this repo

        git clone git@github.com:BadPixel89/labstats-definition-builder.git

Open this in your IDE of choice and insert the URLs you need. Optionally wirte a name parsing function

Navigate to the folder containing the repo (specifically containing main.go) with a terminal and build it

        go build

This should automatically download the required dependencies, if this fails, you may need to install the three libraries this script depends on in order to build it:

        go get github.com/xuri/excelize/v2
        go get github.com/go-ldap/ldap/v3
        go get golang.org/x/term
        go get github.com/BadPixel89/colourtext

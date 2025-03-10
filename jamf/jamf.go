package jamf

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/BadPixel89/colourtext"
	"golang.org/x/term"
)

// This structure is required due to the way the jamf API responds to group requests.
type GroupRoot struct {
	Group ComputerGroup `json:"computer_group"`
}

// The name of a computer group should be it's ID in jamfcloud
type ComputerGroup struct {
	Name      string
	Computers []Computer
}

// computers have other attributes but we only need the name for this
type Computer struct {
	Name string
}

type JamfToken struct {
	Token  string
	Exipry time.Time
}

// TODO
// replace with your subdomain
const (
	jamfIDURL   = "https://<siteurl>.jamfcloud.com/JSSResource/computergroups/id/"
	jamfAuthURL = "https://<siteurl>.jamfcloud.com/api/v1/auth/token"
)

var AuthToken JamfToken = JamfToken{}

// Gets a token from the jamf api using basic auth. This token is cached internally to this package,
// a new token is generated each time auth is called - tokens expire after 30 mins
func JamfAuth() (err error) {
	httpClient := http.Client{}
	req, err := http.NewRequest("POST", jamfAuthURL, nil)
	jamfuser, err := getJamfCredentials("username: ", "enter password: ")
	req.Header.Add("Authorization", "Basic "+jamfuser)
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		colourtext.PrintError(err.Error())
		return err
	}
	defer res.Body.Close()
	//	parse the response into an array of structs and return it
	var apiResponse JamfToken
	err = json.NewDecoder(res.Body).Decode(&apiResponse)
	if err != nil {
		colourtext.PrintError(err.Error())
		return err
	}
	colourtext.PrintSuccess("jamf bearer token received")
	AuthToken = apiResponse
	return err
}

// Takes a list of strings representing group IDs and returns the relevant computer group info from jamfcloud
func SearchJamfByList(Computergroups []string) []ComputerGroup {
	ComputergroupResults := make([]ComputerGroup, 0)

	httpClient := http.Client{}
	for _, group := range Computergroups {
		jamfComputerGroup, err := searchJamfSingle(httpClient, group)
		if err != nil {
			colourtext.PrintError(err.Error())
			continue
		}
		if len(jamfComputerGroup.Computers) == 0 {
			colourtext.PrintWarn("empty computer group found, skipping: " + jamfComputerGroup.Name)
			continue
		}
		ComputergroupResults = append(ComputergroupResults, jamfComputerGroup)
	}

	return ComputergroupResults
}

// Makes an API call that returns the group info of the group with the ID specified in the computergroup parameter.
func searchJamfSingle(client http.Client, computergroup string) (ComputerGroup, error) {
	req, err := http.NewRequest("GET", jamfIDURL+computergroup, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+AuthToken.Token)
	if err != nil {
		colourtext.PrintWarn("failed to create web request for Computergroup: " + computergroup)
		return ComputerGroup{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		colourtext.PrintWarn("Computergroup not found: " + computergroup)
		return ComputerGroup{}, err
	}
	defer res.Body.Close()

	var apiResponse GroupRoot
	err = json.NewDecoder(res.Body).Decode(&apiResponse)
	if err != nil {
		colourtext.PrintWarn("failed to parse response for group: " + computergroup)
		return ComputerGroup{}, err
	}

	colourtext.PrintSuccess(apiResponse.Group.Name)
	return apiResponse.Group, nil
}

// prompt the user for username and password with the given prompt text
func getJamfCredentials(usermsg string, passmsg string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	colourtext.PrintColour(colourtext.Cyan, usermsg)
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	colourtext.PrintColour(colourtext.Cyan, passmsg)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	userstring := strings.TrimSpace(username) + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(userstring)), err
}

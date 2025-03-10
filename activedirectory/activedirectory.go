package activedirectory

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/BadPixel89/colourtext"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/term"
)

type OU struct {
	Name      string
	Computers []string
}

// TODO
// replace these with the correct info for your organisation
// BaseDNAffix is the search root, all nodes you want to pull must be beneath it
const (
	ldapServer  = "ldap://ldap.url.co.uk:389"
	BaseDNAffix = ",OU=Category,OU=Workstations,OU=ROOT,DC=Ldap,DC=url,DC=co,DC=uk"
)

// connect to active directory and return the connection
func ConnectAndBindAD() (*ldap.Conn, error) {
	l, err := ldap.DialURL(ldapServer)
	if err != nil {
		return nil, err
	}
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	username, pass, err := getADCredentials("Active Directory account username@domain.com", "password: ")
	if err != nil {
		return nil, err
	}
	err = l.Bind(username, pass)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// use the provided connection to search for a list of OU names. Errors will be logged to console.
func SearchADbyList(conn *ldap.Conn, ouNames []string) []OU {
	ouResults := make([]OU, 0)
	for _, ou := range ouNames {
		result, err := SearchADSingle(conn, ou)
		if err != nil {
			colourtext.PrintFail("ou: '" + ou + "' not found")
			colourtext.PrintError(err.Error())
			continue
		}
		colourtext.PrintSuccess(ou)
		resultStruct := OU{
			Name: ou,
		}
		for _, obj := range result.Entries {
			//	pulls the last item in the name string, split on comma, then splits on equals to get only the name
			resultStruct.Computers = append(resultStruct.Computers, fmt.Sprint(strings.Split(strings.Split(obj.DN, ",")[0], "=")[1]))
		}
		ouResults = append(ouResults, resultStruct)
	}
	return ouResults
}

// search for all computers in a single OU using the provided connection as search root
func SearchADSingle(l *ldap.Conn, ouDn string) (*ldap.SearchResult, error) {

	filter := "(&(objectClass=computer))"

	searchReq := ldap.NewSearchRequest(
		"OU="+ouDn+BaseDNAffix,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"name"},
		nil,
	)
	result, err := l.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("anonymous Bind Search Error: %s", err)
	}
	if len(result.Entries) > 0 {
		return result, nil
	} else {
		return nil, fmt.Errorf("couldn't fetch anonymous bind search entries: %s", err)
	}
}

// prompt the user for username and password with the given prompt text
func getADCredentials(usermsg string, passmsg string) (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	colourtext.PrintColour(colourtext.Cyan, usermsg)
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	colourtext.PrintColour(colourtext.Cyan, passmsg)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// print an ou to console
func PrintOU(ou OU) {
	fmt.Println(ou.Name)
	for _, comp := range ou.Computers {
		fmt.Println(comp)
	}
}

package main

import (
	"gopkg.in/ldap.v2"
	"os"
	"io/ioutil"
	"strings"
	"fmt"
	"net/url"
)

func main() {
	binddn := ""
	bindpw := ""
	uri := ""
	base := ""
	uid := os.Args[1]

	conf, err := ioutil.ReadFile("/etc/nslcd.conf")
	if err != nil {
		os.Exit(1)
	}
	lines := strings.Split(string(conf), "\n")

	for _, line := range lines {
		words := strings.Split(line, " ")
		if len(words) > 1 {
			key := words[0]
			if key == "uri" {
				uri = words[1]
			} else if key == "base" {
				base = words[1]
			} else if key == "binddn" {
				binddn = words[1]
			} else if key == "bindpw" {
				bindpw = words[1]
			}
		}
	}

	ldapUrl, err := url.Parse(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ldapHost := ldapUrl.Host
	if !strings.Contains(ldapHost, ":") {
		ldapHost += ":389"
	}

	ldapConn, err := ldap.Dial("tcp", ldapHost)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer ldapConn.Close()

	if binddn != "" {
		err = ldapConn.Bind(binddn, bindpw)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	}

	searchRequest := ldap.NewSearchRequest(base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(&(objectClass=posixAccount)(uid=" + uid + "))", []string{"sshPublicKey"}, nil)

	s, err := ldapConn.Search(searchRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	for _, entry := range s.Entries {
		keys := entry.GetAttributeValues("sshPublicKey")
		for _, k := range keys {
			fmt.Println(k)
		}
	}

}
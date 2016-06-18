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
	filter := "(objectClass=posixAccount)"

	conf, err := ioutil.ReadFile("/etc/nslcd.conf")
	if err != nil {
		os.Exit(1)
	}
	lines := strings.Split(string(conf), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			key := fields[0]
			value := strings.Join(fields[1:], " ")
			
			if key == "uri" {
				uri = value
			} else if key == "base" {
				base = value
			} else if key == "binddn" {
				binddn = value
			} else if key == "bindpw" {
				bindpw = value
			}
			
			if len(fields) > 2 && key == "filter" {
				key2 := fields[1]
				value = strings.Join(fields[2:], " ")
				
				if key2 == "passwd" {
					filter = value
				}
			}
		}
	}

	ldapURL, err := url.Parse(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ldapHost := ldapURL.Host
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

	searchRequest := ldap.NewSearchRequest(base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(&" + filter + "(uid=" + uid + "))", []string{"sshPublicKey"}, nil)

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
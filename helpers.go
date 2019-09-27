package httputils

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"gopkg.in/xmlpath.v2"
)

func mapToJSON(m map[string]string) string {
	str, err := json.Marshal(m)

	if err != nil {
		log.Println(err)
		return ""
	}
	return string(str)
}

func JSONToMap(s string) map[string]string {
	var m map[string]string
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		log.Println(err)
	}
	return m
}

func FindXPath(html string, xpath string) (string, error) {
	reader := strings.NewReader(html)
	xmlroot, err := xmlpath.ParseHTML(reader)

	if err != nil {
		log.Println(err)
	}

	// correios table with the package information
	// example `//table[contains(@class, 'listEvent')]`
	path := xmlpath.MustCompile(xpath)
	if value, ok := path.String(xmlroot); ok {
		log.Println("Found:", value)
		return value, nil
	}
	return "", errors.New("No matches with xpath: " + xpath)
}

package client

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

// GetDigest extract digest from key file
func GetDigest(keyFile string) (string, error) {
	xmlFile, err := os.Open(keyFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return "", err
	}

	defer xmlFile.Close()

	type Auth struct {
		Login     string `xml:"login,attr"`
		Realm     string `xml:"realm,attr"`
		MD5       string `xml:"md5,attr"`
		Transport string `xml:"transport,attr"`
		IP        string `xml:"ncc_ip,attr"`
		Port      int    `xml:"ncc_port,attr"`
	}

	type AuthFile struct {
		Auth Auth
	}

	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Println("Error reading key file:", err)
		return "", err
	}

	var authFile AuthFile
	xml.Unmarshal(bytes, &authFile)

	return authFile.Auth.MD5, nil
}

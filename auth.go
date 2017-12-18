package npilib

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

//Auth key data
type auth struct {
	Login     string `xml:"login,attr"`
	Realm     string `xml:"realm,attr"`
	MD5       string `xml:"md5,attr"`
	Transport string `xml:"transport,attr"`
	IP        string `xml:"ncc_ip,attr"`
	Port      int    `xml:"ncc_port,attr"`
}

type authFile struct {
	Auth *auth
}

// GetAuthData extract digest from key file
func getAuthData(keyFile string) (*auth, error) {
	xmlFile, err := os.Open(keyFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return &auth{}, err
	}

	defer xmlFile.Close()

	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Println("Error reading key file:", err)
		return &auth{}, err
	}

	var authFile authFile
	xml.Unmarshal(bytes, &authFile)

	return authFile.Auth, nil
}

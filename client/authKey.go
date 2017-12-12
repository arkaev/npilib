package client

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

//Auth key data
type Auth struct {
	Login     string `xml:"login,attr"`
	Realm     string `xml:"realm,attr"`
	MD5       string `xml:"md5,attr"`
	Transport string `xml:"transport,attr"`
	IP        string `xml:"ncc_ip,attr"`
	Port      int    `xml:"ncc_port,attr"`
}

type authFile struct {
	Auth *Auth
}

// GetAuthData extract digest from key file
func GetAuthData(keyFile string) (*Auth, error) {
	xmlFile, err := os.Open(keyFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return &Auth{}, err
	}

	defer xmlFile.Close()

	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Println("Error reading key file:", err)
		return &Auth{}, err
	}

	var authFile authFile
	xml.Unmarshal(bytes, &authFile)

	return authFile.Auth, nil
}

package utility

import (
	"io/ioutil"
	"os"
	"strings"
)

//ReadFile returns the file content of the passed fileName.
func ReadFile(fileName string) (*string, error) {
	b, err := ioutil.ReadFile(fileName) // just pass the file name
	if err != nil {
		return nil, err
	}
	yaml := string(b)
	return &yaml, nil
}

//GetUserData returns the file content of the passed fileName and replace template variables.
func GetUserData(fileName string) (*string, error) {
	data, err := ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	yaml := parseUserData(*data)
	return &yaml, nil
}

func parseUserData(content string) string {
	clientPublicKey, ok := os.LookupEnv("CLIENT_PUBLICKEY")
	var result string
	if ok == true {
		result = strings.ReplaceAll(content, "{{ CLIENT_PUBLICKEY }}", clientPublicKey)
	} else {
		result = strings.ReplaceAll(content, "{{ CLIENT_PUBLICKEY }}", "")
	}

	metadataURL, ok2 := os.LookupEnv("METADATA_URL")
	if ok2 == true {
		result = strings.ReplaceAll(result, "{{ METADATA_URL }}", metadataURL)
	} else {
		result = strings.ReplaceAll(result, "{{ METADATA_URL }}", "")
	}
	return result
}

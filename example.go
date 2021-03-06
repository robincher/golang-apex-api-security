// Example App to show how to invoke the library and pass in the necessary arguments
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/robincher/golang-apex-api-security/apexsigner"
)

func main() {

	var authParam apexsigner.APIParam
	// Decoding JSON data into Go Values
	var qs map[string]interface{}
	byt := []byte(`{"name":"user1","age":"20"}`)

	if err := json.Unmarshal(byt, &qs); err != nil {
		panic(err)
	}

	// Example Params for HMAC256 Signature
	authParam.Realm = "http://localhost/testapp"
	authParam.AppID = "TestApp"
	authParam.AuthPrefix = "Apex_l2_eg"
	authParam.Secret = "fakesecret"
	authParam.InvokeURL = "http://example.io/api/resource"
	authParam.SignatureURL = "http://example.io/api/resource"
	authParam.HTTPMethod = "GET"
	authParam.QueryString = qs

	result, err := apexsigner.GetAuthorizationHeader(authParam)
	if err != nil {
		fmt.Printf("Error countered")
	}
	fmt.Println("Printing HMAC256 Auth Token")
	fmt.Printf(result + "\n")
	fmt.Println("Sending Mock HTTP Request")
	sendRequest(result)
}

func sendRequest(result string) {
	client := &http.Client{}

	// Please change the URL to your own testing env
	url := "https://httpbin.org/get"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error Init HTTP Request:" + err.Error())
	}

	req.Header.Add("Content-Type", "*")
	req.Header.Add("Authorization", result)
	fmt.Println("Printing Http Request")
	fmt.Println(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error countered when making HTTP Request:" + err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Printing Response")
	fmt.Printf(string(body))

}

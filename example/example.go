// Example App to show how to invoke the library and pass in the necessary arguments
package main

import (
	"encoding/json"
	"fmt"

	"github.com/robincher/golang-apex-api-security"
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
	fmt.Printf(result)
}

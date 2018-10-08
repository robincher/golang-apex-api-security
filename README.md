# Golang HTTP Signature Signer for APEX

A golang http signature library for APEX. It main purpose is to provide a quick utility that generates HTTP security headers for authenticating with secured APEX endpoints

**Security Standards**

1. APEX L1 - HMAC256 HTTP Signature
2. APEX L2 - RSA256 HTTP Signature

*Currently still in designing and testing phase*

## Table of Contents
- [Getting Started](#getting-started)
    * [Installation](#installation)
    * [Walkthrough](#walkthrough)
- [Request Parameters](#request-parameters)
    * [Mandatory Parameters](#madatory-parameters)
    * [Optional Parameters](#optional-parameters)
- [Contributing](#contributing)
- [Release](#release)
- [License](#license)
- [References](#references)

## Getting Started

### Installation

```
go get -u github.com/GovTechSG/golang-apex-api-security
```

If you have errors about package not found , please verify your GOPATH and ensure your project is residing at the right directory. For more information about setting GOPATH, go [here](https://github.com/golang/go/wiki/SettingGOPATH)

### Walkthrough

#### Defining the Request Parameters as a struct

```
// APIParam -  Defining the request struct
type APIParam struct {
	Realm        string `json:"realm"`
	AppID        string `json:"appId"`
	AuthPrefix   string `json:"authPrefix"`
	Secret       string `json:"secret"`
	InvokeURL    string `json:"invokeUrl"`
	SignatureURL string `json:"signatureUrl"`

	HTTPMethod string `json:"httpMethod"`
	Signature  string `json:"signature"`

	PrivateCertFileName string `json:"privateCertFileName"`
	Passphrase          string
	SignatureMethod     string `json:"signatureMethod"`
	Nonce               string `json:"nonce"`
	Timestamp           string `json:"timestamp"`
	Version             string `json:"version"`

	QueryString map[string]interface{} `json:"queryString"`
	FormData    map[string]interface{} `json:"formData"`
}
```

#### Invoking the helper utility

Please appened the values based on your own settings and requirements.

```
requestOpts := APIParam{
		Realm:    "https://portal.example.com,
		AppID: "AppID123",
		AuthPrefix: "apex_lg_l1",
        Secret: "Secret",
        InvokeURL: "https://something.api.gov.sg/v1/resource",
        SignatureURL: ""https://something.api.i.gov.sg/v1/resource""
}

apexAuthHeader, err := getSignatureToken(requestOpts)
```

#### Appending to the HTTP header
```
import net/http example

var req *http.Request
req.Header.Add("Authorization", apexAuthHeader)
```

### Request Parameters

This section describes each of the parameters and it purpose.

#### Mandatory Parameters

1. `AppID`

Apex App ID. The App needs to be approved and activated by the API provider. This value can be obtained from the gateway portal.

2. `AuthPrefix`

API gateway-specific authorization scheme for a specific gateway zone. Takes 1 of 4 possible values.
 
```
var authPrefix = 'Apex_l1_ig'; 
// or
var authPrefix = 'Apex_l1_eg';
// or
var authPrefix = 'Apex_l2_ig';
// or
var authPrefix = 'Apex_l2_eg';
```

3. `HTTPMethod`

The standard HTTP method required for the target resource

4. `InvokeURL`
The full API endpoint path that you will be invoking , for example https://my-apex-api.api.gov.sg/api/my/specific/data

**IMPORTANT NOTE**  
Must be the endpoint URL as served from the Apex gateway, from the domain api.gov.sg. This may differ from the actual HTTP endpoint that you are calling, for example if it were behind a proxy with a different URL.**

5. `SignatureURL`

The API endpoint that is only internally recognised. For examples :
a) https://my-apex-api.i.api.gov.sg/api/my/specific/data 
b) https://my-apex-api.e.api.gov.sg/api/my/specific/data 

We have plans to remove this inconsistency in the near future, so for now do indicate this parameter into the request parameters, else the signature verification will failed.

5a. `Secret`  - For APEX L1 Security 

If the API you are accessing is secured with an APEX L1 policy, you need to provide the generated App secret that corresponds to the `AppID` provided.

5b. `PrivateCertFileName` and `Passphrase`  - For APEX L2 Security 

If the API you are access is secured with an APEX L2 policy, you need to provide the private key and passphrase corresponding to the public key uploaded for `AppID`.

```
var keyPath = "/path/to/my/private.key"
var passphrase = "somepassword"
```

#### Optional Parameters

1.  `Realm`

An identifier for the caller, this can be set to any value. (eg https://portal.example.com)


2. `QueryString` / `FormData`

**IMPORTANT NOTE**  If you pass in the params in QueryString or FormData, please **remove** the query parameters from both SignatureURL and InvokeURL parameter

2b. QueryString

Object representation of URL query parameters, for the API.

For example, the API endpoint is https://example.com/v1/api?key=value , then you have you pass in the params in this manner below :

```
var queryString [][]string
queryString = append(queryString, []string{"key", "value"})

```

2b. `FormData`

Object representation of form data (x-www-form-urlencoded) passed during HTTP POST / HTTP PUT requests

```
var formData [][]string
formData = append(formData, []string{"key", "value"})

```

4. `Nonce`
An arbitrary string, needs to be different after each successful API call. Defaults to 32 bytes random value encoded in base64.

5. `Timestamp`
A unix timestamp. Defaults to the current unix timestamp.

## Contributing
For more information about contributing, and raising PRs or issues, see [CONTRIBUTING.md](https://github.com/GovTechSG/golang-apex-api-security/blob/master/.github/CONTRIBUTING.md).

## Release
See [CHANGELOG.md](CHANGELOG.md).

## License
Licensed under the [MIT LICENSE ](https://github.com/GovTechSG/golang-apex-api-security/blob/master/LICENSE)

## References
+ [RSA and HMAC Request Signing Standard for HTTP Messages](http://tools.ietf.org/html/draft-cavage-http-signatures-09)
+ [99Designs Open-sourced HTTP Signature Golang Library](https://github.com/99designs/httpsignatures-go)
+ [Handling RSA in Go](https://golang.org/src/crypto/rsa/example_test.go?m=text)

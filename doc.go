// Copyright 2018 Government Technology of Singapore , Government Digital Services.
// All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package apexsigner implements utility that conveniently generate a HTTP Signature.

The APEX security standards supported are

1. APEX L1 - HMAC256

2. APEX L2 - RSA256

The name apexsigner stands for "APEX HTTP Signature Signer". Similar to the HTTP signature (HMAC & RSA) standard
defined in IETF, apexsigner generates a HTTP signature with a given set of parameters specified
by the API Gateway. The API Gateway will verified the HTTP signature against the corresponding public key uploaded.

Let's start by creating a request object(struct) for apexsigner

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

Preparing the request options

	For APEX L1 Signature

	The key parameters are the AppID and Secret, which you can retrieved from APEX Community Manager of
	your respective tenant. Additionally, please note the difference between InvokeURL and SignatureURL , which
	is an internally regconised URL. This inconsistency will be removed in the near future.

	requestOpts := APIParam{
			Realm:    "https://portal.example.com,
			AppID: "AppID123",
			AuthPrefix: "apex_lg_l1",
			Secret: "Secret",
			InvokeURL: "https://something.api.gov.sg/v1/resource",
			SignatureURL: ""https://something.api.i.gov.sg/v1/resource""
	}

	For APEX L2 Signature

	The key parameters are the AppID and PrivateCertFileName, which is the file path where you stored your
	private key. Please remember to upload the corresponding public key to the App created in APEX Community
	Manager, and retrieve its AppID.Same as above, please note the difference between InvokeURL and SignatureURL , which
	is an internally regconised URL. This inconsistency will be removed in the near future.

	requestOpts := APIParam{
			Realm:    "https://portal.example.com,
			AppID: "AppID123",
			AuthPrefix: "apex_lg_l2",
			Secret: "Secret",
			PrivateCertFileName: "/path/secret.pm",
			Passphrase: "password"
			InvokeURL: "https://something.api.gov.sg/v1/resource",
			SignatureURL: ""https://something.api.i.gov.sg/v1/resource""
	}

Using the apexsigner library

	apexAuthHeader, err := getSignatureToken(requestOpts)

After getting the Signature Token, just append it to the Http Header's Authorization Scheme

	import net/http example

	var req *http.Request
	req.Header.Add("Authorization", apexAuthHeader)

*/
package apexsigner

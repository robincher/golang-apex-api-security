package apexsigner

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func getSignatureBaseString(baseProps APIParam) (string, error) {
	siteURL, err := url.Parse(baseProps.SignatureURL)
	if err != nil {
		return "", err
	}

	if siteURL.Scheme != "http" && siteURL.Scheme != "https" {
		return "", fmt.Errorf("not supported URL scheme(%s), only http and https allow", siteURL.Scheme)
	}

	// Remove port from url
	signatureURL := fmt.Sprintf("%v://%v%v", siteURL.Scheme, siteURL.Hostname(), siteURL.Path)

	defaultParams, err := getDefaultParam(baseProps)

	// Extract queryString from baseProps
	queryStringArray := parseParam(baseProps.QueryString)
	defaultParams = append(defaultParams, queryStringArray...)

	// Transfer queryString from url for sorting
	query := siteURL.Query()
	siteURL.RawQuery = query.Encode()
	queryMap, _ := url.ParseQuery(siteURL.RawQuery)
	for key, value := range queryMap {
		for _, u := range value {
			defaultParams = append(defaultParams, []string{key, u})
		}
	}

	formDataArray := parseParam(baseProps.FormData)
	defaultParams = append(defaultParams, formDataArray...)

	paramString := ""
	sort.Slice(defaultParams[:], func(i, j int) bool {
		for x := range defaultParams[i] {
			if defaultParams[i][x] == defaultParams[j][x] {
				continue
			}
			return defaultParams[i][x] < defaultParams[j][x]
		}
		return false
	})

	paramString = ArrayNameValuePair{defaultParams}.Stringify()

	sigBaseString := strings.ToUpper(baseProps.HTTPMethod) + "&" + signatureURL + paramString

	return sigBaseString, nil
}

func getSignatureToken(reqProps *APIParam) (string, error) {

	if reqProps.AuthPrefix == "" || reqProps.AppID == "" || reqProps.SignatureURL == "" || reqProps.HTTPMethod == "" {
		return "", fmt.Errorf("one or more required parameters are missing")
	}

	authPrefix := strings.ToLower(reqProps.AuthPrefix)

	if reqProps.Secret == "" {
		reqProps.SignatureMethod = "SHA256withRSA"
	} else {
		reqProps.SignatureMethod = "HMACSHA256"
	}

	if reqProps.Nonce == "" {
		reqProps.Nonce, _ = generateNonce()
	}

	if reqProps.Timestamp == "" {
		reqProps.Timestamp = strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
	}

	baseString, err := getSignatureBaseString(*reqProps)
	if err != nil {
		return "", err
	}

	if reqProps.SignatureMethod == "HMACSHA256" {
		reqProps.Signature, err = getHMACSignature(baseString, reqProps.Secret)
		if err != nil {
			return "", err
		}
	} else {
		privateKey, err := parsePrivateKeyFromPEM(reqProps.PrivateCertFileName, reqProps.Passphrase)
		if err != nil {
			return "", err
		}

		reqProps.Signature, err = getRSASignature(baseString, privateKey)
		if err != nil {
			return "", err
		}
	}

	signatureToken := strings.Title(authPrefix) + " realm=\"" + reqProps.Realm + "\""
	defaultParams, _ := getDefaultParam(*reqProps)
	defaultParams = append(defaultParams, []string{authPrefix + "_signature", reqProps.Signature})

	for _, value := range defaultParams {
		signatureToken = signatureToken + ", " + value[0] + "=\"" + value[1] + "\""
	}

	return signatureToken, nil
}

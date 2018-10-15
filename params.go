package apexsigner

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//APIParam -  API Request parameters struct
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

// ArrayNameValuePair ()
type ArrayNameValuePair struct {
	nameValue [][]string
}

func generateNonce() (string, error) {
	nonceSize := 14

	nonce := make([]byte, nonceSize, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(nonce), nil
}

func getDefaultParam(baseProps APIParam) ([][]string, error) {
	var paramsArray [][]string

	lowerAuthPrefix := strings.ToLower(baseProps.AuthPrefix)
	prefixedAppID := lowerAuthPrefix + "_app_id"
	prefixedNonce := lowerAuthPrefix + "_nonce"
	prefixedSignatureMethod := lowerAuthPrefix + "_signature_method"
	prefixedTimestamp := lowerAuthPrefix + "_timestamp"
	prefixedVersion := lowerAuthPrefix + "_version"

	paramsArray = append(paramsArray, []string{prefixedAppID, baseProps.AppID})

	if baseProps.Nonce == "" {
		nonce, _ := generateNonce()
		paramsArray = append(paramsArray, []string{prefixedNonce, nonce})
	} else {
		paramsArray = append(paramsArray, []string{prefixedNonce, baseProps.Nonce})
	}

	if baseProps.SignatureMethod == "" {
		var signatureMethod string

		if baseProps.Secret == "" {
			signatureMethod = "SHA256withRSA"
		} else {
			signatureMethod = "HMACSHA256"
		}
		paramsArray = append(paramsArray, []string{prefixedSignatureMethod, signatureMethod})
	} else {
		paramsArray = append(paramsArray, []string{prefixedSignatureMethod, baseProps.SignatureMethod})
	}
	if baseProps.Timestamp == "" {
		paramsArray = append(paramsArray, []string{prefixedTimestamp, strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]})
	} else {
		paramsArray = append(paramsArray, []string{prefixedTimestamp, baseProps.Timestamp})
	}

	if baseProps.Version == "" {
		paramsArray = append(paramsArray, []string{prefixedVersion, "1.0"})
	} else {
		paramsArray = append(paramsArray, []string{prefixedVersion, baseProps.Version})
	}

	return paramsArray, nil
}

// Search ()
func (param ArrayNameValuePair) searchFirst(name string) string {
	for _, value := range param.nameValue {
		if value[0] == name {
			return value[1]
		}
	}

	return ""
}

// Stringify - Stringify request parameters
func (param ArrayNameValuePair) Stringify() string {
	paramString := ""

	for _, value := range param.nameValue {
		if value[1] == "" {
			paramString = fmt.Sprintf("%v&%v", paramString, value[0])
		} else {
			paramString = fmt.Sprintf("%v&%v=%v", paramString, value[0], value[1])
		}
	}

	return paramString
}

// parseParam - Parse query parameters and return as array
func parseParam(params map[string]interface{}) [][]string {
	var paramArray [][]string

	for key, value := range params {
		switch v := value.(type) {
		case string:
			paramArray = append(paramArray, []string{key, v})
		case bool:
			paramArray = append(paramArray, []string{key, strconv.FormatBool(v)})
		case float64:
			paramArray = append(paramArray, []string{key, strconv.FormatFloat(v, 'f', -1, 64)})
		case []interface{}:
			paramInterfaceArray := parseArray(key, v)
			paramArray = append(paramArray, paramInterfaceArray...)
		default: // support string, bool, float64 and array only, all other datatype will be ignore
			paramArray = append(paramArray, []string{key, ""})
			//return nil, fmt.Errorf("paramsStringify - params fields data type for '%s' of type '%s' not supported", key, reflect.TypeOf(v).Kind())
		}
	}

	return paramArray
}

// parseArray - Parse array
func parseArray(key string, params []interface{}) [][]string {
	var paramArray [][]string

	for _, value := range params {
		switch v := value.(type) {
		case string:
			paramArray = append(paramArray, []string{key, v})
		case bool:
			paramArray = append(paramArray, []string{key, strconv.FormatBool(v)})
		case float64:
			paramArray = append(paramArray, []string{key, strconv.FormatFloat(v, 'f', -1, 64)})
		default: // support string, bool, float64 only, all other datatype will be ignore
			paramArray = append(paramArray, []string{key, ""})
			//return nil, fmt.Errorf("paramArrayStringify - params fields data type for '%s' of type '%s' not supported", key, reflect.TypeOf(v).Kind())
		}
	}
	return paramArray
}

// GetSignatureToken - Public function to get Apex Signature authorization token to be append on HTTP header
func GetSignatureToken(reqProps APIParam) (string, error) {
	return getSignatureToken(&reqProps)
}

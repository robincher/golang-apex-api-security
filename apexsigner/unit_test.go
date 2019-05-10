package apexsigner

import (
	"crypto/rsa"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var BASEPATH = "../helpers/"
var STUBPATH = fmt.Sprintf("%s/testdata/", BASEPATH)

//var BASEPATH = fmt.Sprintf("%s/src/github.com/GovTechSG/test-suites-apex-api-security/", os.Getenv("GOPATH"))

func Test_defaultParams(t *testing.T) {
	executeTest(t, STUBPATH+"params.json", func(param *TestParam) (string, error) {
		result, err := getDefaultParam(param.APIParam)

		// timestamp value not set in input param, update the expected result after getDefaultParam set the value
		if param.APIParam.Timestamp == "" {
			param.setExpectedResult(fmt.Sprintf(param.getExpectedResult(), (ArrayNameValuePair{nameValue: result}).searchFirst(strings.ToLower(param.APIParam.AuthPrefix)+"_timestamp")))
		}

		// nonce value not set in input param, update the expected result after getDefaultParam set the value
		if param.APIParam.Nonce == "" {
			param.setExpectedResult(fmt.Sprintf(param.getExpectedResult(), (ArrayNameValuePair{nameValue: result}).searchFirst(strings.ToLower(param.APIParam.AuthPrefix)+"_nonce")))
		}

		paramString := (ArrayNameValuePair{nameValue: result}).Stringify()

		return paramString, err
	})
}

func Test_l1Signature(t *testing.T) {
	executeTest(t, STUBPATH+"hmac_sig.json", func(param *TestParam) (string, error) {
		return getHMACSignature(param.Message, param.APIParam.Secret)
	})
}

func Test_L2Signature(t *testing.T) {
	executeTest(t, STUBPATH+"rsa_sig.json", func(param *TestParam) (string, error) {
		result := ""

		privateKey, err := parsePrivateKeyFromPEM(BASEPATH+param.APIParam.PrivateCertFileName, param.APIParam.Passphrase)
		if err == nil {
			result, err = getRSASignature(param.Message, privateKey)
		}

		return result, err
	})
}

func Test_getSignatureBaseString(t *testing.T) {
	executeTest(t, STUBPATH+"basestring.json", func(param *TestParam) (string, error) {
		return getSignatureBaseString(param.APIParam)
	})
}

func Test_getSignatureToken(t *testing.T) {
	executeTest(t, STUBPATH+"sig_token.json", func(param *TestParam) (string, error) {
		dynamicTimestamp := false
		if param.APIParam.Timestamp == "" {
			dynamicTimestamp = true
		}
		dynamicNonce := false
		if param.APIParam.Nonce == "" {
			dynamicNonce = true
		}

		// reset the cert folder...
		param.APIParam.PrivateCertFileName = BASEPATH + param.APIParam.PrivateCertFileName

		result, err := getSignatureToken(&param.APIParam)

		if err == nil {
			if dynamicTimestamp && dynamicNonce {
				param.setExpectedResult(fmt.Sprintf(param.getExpectedResult(), param.APIParam.Nonce, param.APIParam.Timestamp, param.APIParam.Signature))
			} else if dynamicTimestamp {
				param.setExpectedResult(fmt.Sprintf(param.getExpectedResult(), param.APIParam.Timestamp, param.APIParam.Signature))
			} else if dynamicNonce {
				param.setExpectedResult(fmt.Sprintf(param.getExpectedResult(), param.APIParam.Nonce, param.APIParam.Signature))
			}
		}
		return result, err
	})
}

func Test_verifyL1Signature(t *testing.T) {
	executeTest(t, STUBPATH+"hmac_sig_verify.json", func(param *TestParam) (string, error) {
		result, err := VerifyHMACSignature(param.Message, param.APIParam.Secret, param.APIParam.Signature)
		return strconv.FormatBool(result), err
	})
}

func Test_verifyL2Signature(t *testing.T) {
	executeTest(t, STUBPATH+"rsa_sig_verify.json", func(param *TestParam) (string, error) {
		//t.Logf(">>> filetype %s -- %s <<<", filepath.Ext(param.PublicCertFileName), param.PublicCertFileName)
		result := false
		var publicKey *rsa.PublicKey
		var err error

		fileExtension := strings.ToLower(filepath.Ext(BASEPATH + param.PublicCertFileName))

		if fileExtension == ".cer" {
			publicKey, err = parsePublicKeyFromCert(BASEPATH + param.PublicCertFileName)
		} else if fileExtension == ".pem" {
			//publicKey, err = getPublicKeyFromPEM(param.PublicCertFileName)
			publicKey, err = parsePublicKeyFromCert(BASEPATH + param.PublicCertFileName)
		} else if fileExtension == ".key" {
			publicKey, err = parsePublicKeyFromKey(BASEPATH + param.PublicCertFileName)
		} else {
			t.Errorf("\nnot supported file tyep::: %s", BASEPATH+param.PublicCertFileName)
		}

		if err == nil {
			result, err = VerifyRSASignature(param.Message, publicKey, param.APIParam.Signature)
		}

		return strconv.FormatBool(result), err
	})
}

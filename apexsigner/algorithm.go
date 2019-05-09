package apexsigner

import (
	"crypto"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// getRSASignature - Get a RSA Signature given the basestring and private key
// It returns the RSA Signature
func getRSASignature(message string, privateKey *rsa.PrivateKey) (string, error) {
	if message == "" || privateKey == nil {
		return "", fmt.Errorf("message and privateKey must not be null or empty")
	}

	base64Signature, err := signPKCS1v15(privateKey, message)
	if err != nil {
		return "", err
	}

	return base64Signature, nil
}

//getHMACSignature - Get a HMAC256 Signature given basestring and app secret
// It returns the HMAC Signature
func getHMACSignature(message string, secret string) (string, error) {
	if message == "" || secret == "" {
		return "", fmt.Errorf("message and secret must not be null or empty")
	}

	messageHMAC := hmac.New(sha256.New, []byte(secret))
	messageHMAC.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(messageHMAC.Sum(nil)), nil
}

//https://golang.org/src/crypto/rsa/example_test.go?m=text
func signPKCS1v15(rsaPrivateKey *rsa.PrivateKey, message string) (base64Signature string, err error) {
	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.

	byteMessage := []byte(message)

	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed. This requires
	// that the hash function be collision resistant. SHA-256 is the
	// least-strong hash function that should be used for this at the time
	// of writing (2016).
	hashed := sha256.Sum256(byteMessage)

	signature, err := rsa.SignPKCS1v15(nil, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	base64Signature = base64.StdEncoding.EncodeToString(signature)

	return base64Signature, nil
}

// VerifyHMACSignature - Verify HMAC256 signature given a basestring , app secret and signature
// It returns the HMAC256 signature
func VerifyHMACSignature(message string, secret string, signature string) (bool, error) {
	newSignature, err := getHMACSignature(message, secret)
	if err != nil {
		return false, err
	}
	return newSignature == signature, nil
}

// VerifyRSASignature - Verify RSA256 signature given a basestring , corresponding public key and signature
// This is utilizing native PKCS 1.5 encryption standard
// It returns a boolean for the verification
func VerifyRSASignature(message string, rsaPublicKey *rsa.PublicKey, base64Signature string) (bool, error) {
	// convert string to []byte
	byteMessage := []byte(message)

	// convert base64 to []byte
	signature, err := base64.StdEncoding.DecodeString(base64Signature)
	if err != nil {
		return false, err
	}

	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed. This requires
	// that the hash function be collision resistant. SHA-256 is the
	// least-strong hash function that should be used for this at the time
	// of writing (2016).
	hashed := sha256.Sum256(byteMessage)

	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, err
	}

	// signature is a valid signature of message from the public key.
	return true, nil
}

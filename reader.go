package apexsigner

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// parsePublicKeyFromKey - Retrieve the publickey from a DER encoded public key (.key extension)
// It returns the public key contents
func parsePublicKeyFromKey(publicKeyFilePath string) (*rsa.PublicKey, error) {
	// Read the private key
	pemData, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

// parsePublicKeyFromCert - Retrieve the publickey from a .pem or .cert file extension
// It returns the public key contents
func parsePublicKeyFromCert(publicKeyFilePath string) (*rsa.PublicKey, error) {
	// Read the private key
	pemData, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %v", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %s", err.Error())
	}

	return cert.PublicKey.(*rsa.PublicKey), nil
}

// parsePrivateKeyFromPEM - Retrieve the private from a .pem  file extension
// It returns the private key contents
func parsePrivateKeyFromPEM(privateKeyFilePath string, passphrase string) (privateKey *rsa.PrivateKey, err error) {
	// Read the private key
	pemData, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	if got, want := block.Type, "RSA PRIVATE KEY"; got == want {
		//log.Fatalf("unknown key type %q, want %q", got, want)
		if passphrase != "" {
			decBlock, err := x509.DecryptPEMBlock(block, []byte(passphrase))
			if err != nil {
				//log.Fatalf("error decrypting pem file: %s", err.Error())
				return nil, err
			}
			block.Bytes = decBlock
		}
		// Decode the RSA private key
		parseResult, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			//log.Fatalf("bad private key: %s", err)
			return nil, err
		}

		privateKey = parseResult
	}

	if got, want := block.Type, "PRIVATE KEY"; got == want {
		// Decode the RSA private key
		parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		// cast to correct type, rsaPrivateKey
		privateKey = parseResult.(*rsa.PrivateKey)
	}

	if privateKey == nil {
		path := privateKeyFilePath
		fileName := filepath.Base(path)
		return nil, fmt.Errorf("Unable to read private key from %v", fileName)
	}

	return privateKey, nil
}

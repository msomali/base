package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)


// PinEncryptionRSA encrypts the pin with the public key
// returns a base64 encoded string of the encrypted pin
func PinEncryptionRSA(pin string, pubKey string) (string, error) {

	decodedBase64, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return "", fmt.Errorf("could not decode pub key to Base64 string: %w", err)
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedBase64)
	if err != nil {
		return "", fmt.Errorf("could not parse encoded public key (encryption key) : %w", err)
	}

	//check if the public key is RSA public key
	publicKey, isRSAPublicKey := publicKeyInterface.(*rsa.PublicKey)
	if !isRSAPublicKey {
		return "", fmt.Errorf("public key parsed is not an RSA public key : %w", err)
	}
	msg := []byte(pin)
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
	if err != nil {
		return "", fmt.Errorf("could not encrypt api key using generated public key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil

}


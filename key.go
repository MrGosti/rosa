package rosa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/user"
)

func savePrivateKey(key *rsa.PrivateKey, filename string) error {

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	err := saveFile(filename, pemdata)
	if err != nil {
		return err
	}
	return nil
}

// StringifyPublicKey create a base64 encoded version of the public key, usefull for Friends saving or sharing
func StringifyPublicKey(key *rsa.PublicKey) string {
	return fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(key.N.Bytes()))
}

// UnStringifyPublicKey create *rsa.PublicKey from a base64 encoded string, assuming that the exposant is 65537 (See wikipedia for further information). Throw an error if the key is not valid
func UnStringifyPublicKey(content string) (*rsa.PublicKey, error) {

	N := big.NewInt(0)
	key, err := base64.StdEncoding.DecodeString(content)
	N = N.SetBytes(key)

	return &rsa.PublicKey{N, 65537}, err
}

func savePublicKey(key *rsa.PublicKey, identifier string, filename string) error {

	err := saveFile(filename, []byte(identifier+" "+StringifyPublicKey(key)+"\n"))
	if err != nil {
		return err
	}
	return nil
}

// LoadPrivateKey open a file and return you the PEM and PKCS1 decoded *rsa.PrivateKey
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	filecontent, err := loadFile(filename)
	if err != nil {
		return nil, err
	}
	key, _ := pem.Decode(filecontent)
	if key == nil {
		return nil, errors.New("The file is not valid")
	}
	privatekey, err := x509.ParsePKCS1PrivateKey(key.Bytes)

	return privatekey, err
}

func isPrivKeyAvailable() bool {
	usr, err := user.Current()
	if err != nil {
		return false
	}

	if _, err := os.Stat(usr.HomeDir + "/.rosa/key.priv"); err == nil {
		return true
	}
	return false
}

func isPubKeyAvailable() bool {

	usr, err := user.Current()
	if err != nil {
		return false
	}

	if _, err := os.Stat(usr.HomeDir + "/.rosa/key.pub"); err == nil {
		return true
	}
	return false
}

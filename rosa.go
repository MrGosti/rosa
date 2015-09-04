package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os/user"
)

func Decrypt(content []byte, privatekey *rsa.PrivateKey) ([]byte, error) {
	md5hash := md5.New()
	label := []byte("")

	decryptedmsg, err := rsa.DecryptOAEP(md5hash, rand.Reader, privatekey, content, label)

	if err != nil {
		return nil, err
	}
	return decryptedmsg, nil
}

func Encrypt(content []byte, publickey *rsa.PublicKey) ([]byte, error) {
	md5hash := md5.New()
	label := []byte("")

	encryptedmsg, err := rsa.EncryptOAEP(md5hash, rand.Reader, publickey, content, label)
	if err != nil {
		return nil, err
	}
	return encryptedmsg, nil
}

func Generate(identifier string, save bool) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	var publickey *rsa.PublicKey
	var privatekey *rsa.PrivateKey

	usr, err := user.Current()
	privatekey, err = rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		return nil, nil, err
	}

	privatekey.Precompute()
	err = privatekey.Validate()

	if err != nil {
		return nil, nil, err
	}

	publickey = &privatekey.PublicKey

	if save == true {
		savePrivateKey(privatekey, usr.HomeDir+"/.rosa/key.priv")
		savePublicKey(publickey, identifier, usr.HomeDir+"/.rosa/key.pub")
	}
	return privatekey, publickey, nil
}

func main() {
	usr, _ := user.Current()
	Generate(usr.Username, true)
	_, err := LoadPrivateKey(usr.HomeDir + "/.rosa/key.priv")
	if err != nil {
		fmt.Println(err)
	}

	LoadFriends(usr.HomeDir + "/.rosa/friend_list")
	fmt.Println(len(FriendList))
	FriendList["8b49905dcce57a634e18e386aa7f6b59"].Remove("fe")
	fmt.Println(len(FriendList))
	fmt.Println(FriendList["8b49905dcce57a634e18e386aa7f6b59"])
	// for i := 0; i < 40; i++ {
	// 	name := fmt.Sprintf("Test%d", i)
	// 	_, pub, _ := Generate(name, false)
	// 	f := &Friend{name, pub}
	// 	f.Registrer(usr.HomeDir + "/.rosa/friend_list")
	// }

	// msg, _ := Encrypt([]byte("Hello world"), publickey)
	// decrypted, _ := Decrypt(msg, wierd)

	// fmt.Printf("%v\n", string(decrypted))
}

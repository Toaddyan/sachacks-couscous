package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

func generateSSHKey() string {
	savePrivateFileTo := "./secrets/id_rsa"
	savePublicFileTo := "./secrets/id_rsa.pub"
	bitSize := 4096

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err.Error())
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	err = writeKeyToFile(privateKeyBytes, savePrivateFileTo)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = writeKeyToFile([]byte(publicKeyBytes), savePublicFileTo)
	if err != nil {
		log.Fatal(err.Error())
	}

	publicKey := string(publicKeyBytes)

	// Removing trailing character from publicKey
	r := []rune(publicKey)
	publicKey = string(r[:len(publicKey)-1])

	return publicKey
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}

// func generateSSHKey() string {

// 	// Generate privateKey
// 	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
// 	if err != nil {
// 		fmt.Printf("Cannot generate RSA key\n")
// 		os.Exit(1)
// 	}

// 	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
// 	privateKeyBlock := &pem.Block{
// 		Type:  "RSA PRIVATE KEY",
// 		Bytes: privateKeyBytes,
// 	}

// 	privatePem, err := os.Create("secrets/private.pem")
// 	if err != nil {
// 		fmt.Printf("error when create private.pem: %s \n", err)
// 		os.Exit(1)
// 	}
// 	err = pem.Encode(privatePem, privateKeyBlock)
// 	if err != nil {
// 		fmt.Printf("error when encode private pem: %s \n", err)
// 		os.Exit(1)
// 	}

// 	publickey := &privatekey.PublicKey
// 	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
// 	if err != nil {
// 		fmt.Printf("error when dumping publickey: %s \n", err)
// 		os.Exit(1)
// 	}
// 	publicKeyBlock := &pem.Block{
// 		Type:  "PUBLIC KEY",
// 		Bytes: publicKeyBytes,
// 	}

// 	publicPem, err := os.Create("secrets/public.pem")
// 	if err != nil {
// 		fmt.Printf("error when create public.pem: %s \n", err)
// 		os.Exit(1)
// 	}
// 	err = pem.Encode(publicPem, publicKeyBlock)
// 	if err != nil {
// 		fmt.Printf("error when encode public pem: %s \n", err)
// 		os.Exit(1)
// 	}
// 	p, err := ioutil.ReadAll(publicPem)
// 	if err != nil {
// 		log.Panic("read failed")
// 	}
// 	return string(p)

// }

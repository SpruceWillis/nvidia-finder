package auth

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// Credentials hold creds and server info for email
type Credentials struct {
	Username, Password, Hostname, Port string
}

// GetPlainAuth parse credentials into a smtp.Auth object
func (creds *Credentials) GetPlainAuth() *smtp.Auth {
	auth := smtp.PlainAuth("", creds.Username, creds.Password, creds.Hostname)
	return &auth
}

// GetCreds read credentials from file or keyboard input, saving if the latter
func GetSmtpCreds() *Credentials {
	credsFile := "credentials.txt"
	creds, err := credsFromFile(credsFile)
	if err != nil {
		// initialize this variable so that creds is not recreated only inside this scope
		var key string
		fmt.Println("credentials not found")
		creds, key = getCredsFromInput()
		saveCreds(creds, credsFile, key)
	}
	return creds
}

func credsFromFile(fileName string) (*Credentials, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("unable to read credentials from file", err)
		return nil, err
	}
	fmt.Println(file)
	// TODO finish this
	return nil, errors.New("not yet implemented")
}

// prompt the user for credentials, hostname, and port from keyboard
func getCredsFromInput() (*Credentials, string) {
	fmt.Println("requesting user credentials")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("username: ")
	username, _ := reader.ReadString('\n')
	username = stripString(username)
	fmt.Println("password: ")
	password, _ := reader.ReadString('\n')
	password = stripString(password)
	fmt.Println("email server host: ")
	hostname, _ := reader.ReadString('\n')
	hostname = stripString(hostname)
	fmt.Println("port: ")
	port, _ := reader.ReadString('\n')
	port = stripString(port)
	fmt.Println("enter encryption key if so desired, leave blank otherwise")
	key, _ := reader.ReadString('\n')
	key = stripString(key)
	creds := Credentials{username, password, hostname, port}
	return &creds, key
}

func saveCreds(creds *Credentials, credsFile, key string) error {
	encryptionEnabled := key != ""
	credentialJSON, err := json.Marshal(&creds)
	if err != nil {
		log.Println("unable to save credentials as JSON, skipping", err)
		return err
	}
	var credBytes []byte
	if encryptionEnabled {
		credBytes, err = encryptCreds(credentialJSON, []byte(key))
		if err != nil {
			log.Println("encryption error while writing credentials", err)
			return err
		}
	} else { // just write in plaintext
		credBytes = credentialJSON
	}
	err = ioutil.WriteFile(credsFile, credBytes, 0644)
	return err
}

// just use AES for
func encryptCreds(credsData, key []byte) ([]byte, error) {
	// initialize block cipher
	c, err := aes.NewCipher(key)
	// typically errors are due to improper key sizes
	if err != nil {
		log.Println("unable to use key", key, "as AES key")
		return nil, err
	}
	// run block cipher in GCM mode, which is apparently faster
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println("unable to use key", key, "as AES key")
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("unable to populate crytographic nonce")
		return nil, err
	}

	return gcm.Seal(nonce, nonce, credsData, nil), nil
}

func stripString(s string) string {
	return strings.Replace(s, "\n", "", -1)
}

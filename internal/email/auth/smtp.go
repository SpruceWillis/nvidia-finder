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

// GetURL return the Host:Port from these credentials
func (creds *Credentials) GetURL() string {
	return fmt.Sprint(creds.Hostname, ":", creds.Port)
}

// GetSmtpCreds read credentials from file or keyboard input, saving if the latter
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
	var key string
	var creds Credentials
	if _, err := os.Stat(fileName); err != nil {
		log.Println("no credentials file found, skipping")
		return nil, err
	}
	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(fileData, &creds)
	if err != nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("unable to read credentials from file, please enter encryption key")
		key, _ = reader.ReadString('\n')
		c, err := decryptCredentials(fileData, key)
		return c, err
	}
	return &creds, nil
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
	fmt.Println("enter aes key if so desired")
	fmt.Println("valid keys are 16, 24, or 32 bytes; you can also enter these numbers to randomly generate a key")
	fmt.Println("leave blank to store credentials in plaintext")
	keyString, _ := reader.ReadString('\n')
	keyString = stripString(keyString)
	key := parseKeyOpt(keyString)
	creds := Credentials{username, password, hostname, port}
	return &creds, key
}

func parseKeyOpt(keyString string) string {
	switch keyString {
	case "16":
		return generateRandomAesKey(16)
	case "24":
		return generateRandomAesKey(24)
	case "32":
		return generateRandomAesKey(32)
	default:
		return keyString
	}
}

func generateRandomAesKey(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		// we don't really expect this to happen, but maybe we ran out of entropy?
		log.Println("error:", err)
	}
	result := string(b)
	log.Println("encrypting credentials with key", result, "please store this somewhere safe")
	return result
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

// TODO use bcrypt package to generate key and encrypt using that
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

func decryptCredentials(encryptedCreds []byte, key string) (*Credentials, error) {
	var creds Credentials
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Println(err)
		return &Credentials{}, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println(err)
		return &Credentials{}, err
	}
	nonceSize := gcm.NonceSize()
	if len(encryptedCreds) < nonceSize {
		err = errors.New("nonce shorter than credentials cyphertext, unable to decrypt")
		log.Println(err)
		return &Credentials{}, err
	}

	nonce, encryptedCreds := encryptedCreds[:nonceSize], encryptedCreds[nonceSize:]
	decryptedCreds, err := gcm.Open(nil, nonce, encryptedCreds, nil)
	if err != nil {
		log.Println(err)
		return &Credentials{}, err
	}
	err = json.Unmarshal(decryptedCreds, &creds)
	if err != nil {
		log.Println(err)
		return &Credentials{}, err
	}
	return &creds, nil
}

func stripString(s string) string {
	return strings.Replace(s, "\n", "", -1)
}

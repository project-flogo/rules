package common

import (
	"crypto/rand"
	"io"
	"time"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/oklog/ulid"
)

//unique id (26 chars)
func GetUniqueId() (string, error) {
	entropy, err := generateRandomBytes(10)
	if err != nil {
		return "", err
	}
	uuid, err := ulid.New(ulid.Timestamp(time.Now().UTC()), entropy)
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) (io.Reader, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return rand.Reader, nil
}

// GetAbsPathForResource gets the absolute resource path off $GOPATH
func GetAbsPathForResource(resourcepath string) string {
	GOPATH := os.Getenv("GOPATH")
	fmt.Printf("GOPATH - [%s]\n", GOPATH)
	paths := strings.Split(GOPATH, ":")
	for _, path := range paths {
		// fmt.Printf("path[%s]\n", path)
		absPath := path + "/" + resourcepath
		_, err := os.Stat(absPath)
		if err == nil {
			return absPath
		}
	}
	return ""
}

// FileToString reads file content to string
func FileToString(fileName string) string {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(dat)
}

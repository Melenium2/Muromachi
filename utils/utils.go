package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
)

// generate random uuid sting like 076aa5bb-7f92-4410-9fd3-23d5916a5796
func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// Generate random sha256 hash string from given str and random
// salt
func Hash(str string, random interface{}) string {
	b := &bytes.Buffer{}

	enc := gob.NewEncoder(b)
	err := enc.Encode(random)
	if err != nil {
		panic(fmt.Sprintf("terrible gob encoded content, %v", err))
	}

	salt := make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		panic(fmt.Sprintf("terrible hash, %v", err))
	}

	_, _ = b.Write(salt)
	_, _ = b.WriteString(str)

	hash := sha256.Sum256(b.Bytes())

	return fmt.Sprintf("%x", hash)
}

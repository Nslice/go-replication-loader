package argsp

import (
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"io"
	"os"

	"github.com/sergeyzalunin/go-replication-loader/logger"
)

const (
	filename       = "data.dat"
	masterPassword = "SUPERSECRETPASSWSUPERSECRETPASSW" // 32
)

// Serialize arguments by marshaling it in json.
// Further this json object encodes by aes and writes into file by using gob package
func Serialize(args *ArgumentOptions, log *logger.Log) {
	encodedArgs, err := encodeArguments(args)
	if err != nil {
		log.Error(err)
		return
	}

	err = writeArgs(encodedArgs)
	if err != nil {
		log.Error(err)
	}
}

// Deserialize arguments
func Deserialize(log *logger.Log) *ArgumentOptions {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &ArgumentOptions{}
	}
	
	encodedArgs, err := readGob()
	if err != nil {
		log.Error(err)
		return &ArgumentOptions{}
	}

	obj, err := decodeArguments(encodedArgs)
	if err != nil {
		log.Error(err)
		return &ArgumentOptions{}
	}

	args := obj.(ArgumentOptions)
	return &args
}

func writeArgs(encodedArgs []byte) error {
	if file, err := os.Create(filename); err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(encodedArgs)
		file.Close()
	}

	return nil
}

func readGob() ([]byte, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err == nil {
		var encodedArgs []byte
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&encodedArgs)
		if err != nil {
			return nil, err
		}
		return encodedArgs, nil
	}

	return nil, err
}

func encodeArguments(obj interface{}) ([]byte, error) {
	text, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("encodeArguments, Marshal args: \n%v", err)
	}

	block, err := aes.NewCipher([]byte(masterPassword))
	if err != nil {
		return nil, fmt.Errorf("encodeArguments, getting new chipher: \n%v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("encodeArguments, getting new GCM: \n%v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("encodeArguments, reading all bytes: \n%v", err)
	}

	result := gcm.Seal(nonce, nonce, text, nil)
	return result, nil
}

func decodeArguments(encodedArgs []byte) (interface{}, error) {
	block, err := aes.NewCipher([]byte(masterPassword))
	if err != nil {
		return nil, fmt.Errorf("dencodeArguments, getting new chipher: \n%v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("dencodeArguments, getting new GCM: \n%v", err)
	}

	
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := encodedArgs[:nonceSize], encodedArgs[nonceSize:]

	obj, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("dencodeArguments, opening reader to encode args: \n%v", err)
	}

	var result ArgumentOptions
	err = json.Unmarshal(obj, &result)
	if err != nil {
		return nil, fmt.Errorf("dencodeArguments, Unmarshal args: \n%v", err)
	}

	return result, nil
}

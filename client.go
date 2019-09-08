package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	encryptURI      = "http://localhost:8080/encrypt"
	decryptURI      = "http://localhost:8080/decrypt"
	applicationType = "application/json"
)

type encryptionRequest struct {
	ID   string
	Data string
}

type encryptionResponse struct {
	Result string
	Key    []byte
}

type EncrypterClient struct {
	id      string
	payload []byte
}

func (e EncrypterClient) Store(id, payload []byte) (aesKey []byte, err error) {
	encryptionReq := encryptionRequest{}
	encryptionReq.ID = string(id)
	encryptionReq.Data = string(payload)

	reqBody, err := json.Marshal(encryptionReq)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(encryptURI, applicationType, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response encryptionResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return response.Key, nil
}

type decryptionRequest struct {
	ID  string
	Key []byte
}

type decryptionResponse struct {
	Result string
	Data   string
}

func (e EncrypterClient) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	decryptionReq := decryptionRequest{}
	decryptionReq.ID = string(id)
	decryptionReq.Key = aesKey

	reqBody, err := json.Marshal(decryptionReq)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(decryptURI, applicationType, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response decryptionResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return []byte(response.Data), nil
}

func main() {
	var clinet Client
	clinet = EncrypterClient{}

	aesKey, err := clinet.Store([]byte("1"), []byte("Yoti"))
	if err != nil {
		fmt.Printf("Error while trying to encrypt: %v", err)
	}

	plainText, err := clinet.Retrieve([]byte("1"), aesKey)
	if err != nil {
		fmt.Printf("Error while trying to decrypt: %v", err)
	}

	fmt.Printf("Response ... Plain Text: %v", string(plainText))
}

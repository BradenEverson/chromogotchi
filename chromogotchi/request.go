package main

import (
	"encoding/json"
	"errors"
)

type RequestObject struct {
	RequestType string `json:"type"`
	Metadata    []byte `json:"data"`
}

func makeRequestObject(reqtype string, data []byte) RequestObject {
	return RequestObject{reqtype, data}
}

func deserializeRequestObject(data []byte) (*RequestObject, error) {
	object := new(RequestObject)
	err := json.Unmarshal(data, object)
	if err != nil {
		return object, nil
	} else {
		return nil, errors.New("Failed to deserialize request")
	}
}

func serializeRequestObject(request *RequestObject) (*string, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("Failed to serialize request")
	} else {
		requestString := string(data)
		return &requestString, nil
	}
}

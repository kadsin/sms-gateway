package dtos

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

func Marshal(data any) (b []byte, err error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return marshaledData, nil
}

func Unmarshal[T any](source []byte) (T, error) {
	var data T
	if err := json.Unmarshal(source, &data); err != nil {
		return data, err
	}

	if err := validator.New().Struct(&data); err != nil {
		return data, err
	}

	return data, nil
}

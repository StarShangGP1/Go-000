package Dao

import (
	"GoProject/Week02/DB"
	"github.com/pkg/errors"
)

func ErrorMessages() (string, error) {
	value, err := DB.ThrowError()
	if err != nil {
		return "", errors.WithMessage(err, "Wrap errors DAO")
	}
	return value, nil
}

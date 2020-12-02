package Service

import (
	"GoProject/Week02/Dao"
	"github.com/pkg/errors"
)

func WrapError() (string, error) {
	value, err := Dao.ErrorMessages()

	if err != nil {
		return "", errors.Wrap(err, "Wrap errors service")
	}
	return value, nil
}

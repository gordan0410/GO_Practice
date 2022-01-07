package tools

import (
	"errors"
	"strconv"
)

func Parse_int(n interface{}) (int, error) {
	s, ok := n.(string)
	if !ok {
		return 0, errors.New("interface convert to string error")
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("string convert to int error")
	}
	return i, nil
}

package dbtable

import (
	"errors"
	"strings"
)

func parseKeyValue(resPath string) (key string, value string, err error) {
	pos := strings.LastIndex(resPath, "/") + 1

	if pos == 0 {
		err = errors.New("key path is in wrong format")
		return
	}

	if pos >= len(resPath) {
		err = errors.New("key or value is not provided")
		return
	}

	splitted := strings.SplitN(resPath[pos:], ":", 2)

	if len(splitted) == 1 {
		value = splitted[0]
	} else {
		key = splitted[0]
		value = splitted[1]
	}

	return
}

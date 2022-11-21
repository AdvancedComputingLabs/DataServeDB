package dbtable

import (
	"errors"
	"strings"
)

//TODO: make extraction function of keyword or function name and return its type.
// Above is better than below?

func isFunction(key string) bool {
	//a function is that starts with a '$' sign, has  a '(' sign before ')', and ends with ')'.
	// e.g. $func_name(arg1, arg2, arg3)
	// e.g. (not supported at the moment) $func_name(arg1, arg2, arg3, $func_name2(arg1, arg2, arg3))

	// check if it is function
	if len(key) < 3 {
		return false
	}

	if key[0] != '$' {
		return false
	}

	if key[len(key)-1] != ')' {
		return false
	}

	// check if it has '('
	pos := strings.Index(key, "(")
	if pos == -1 {
		return false
	}

	return true
}

func extractFunctionName(key string) (string, error) {
	//a function is that starts with a '$' sign, has  a '(' sign before ')', and ends with ')'.
	// e.g. $func_name(arg1, arg2, arg3)
	// e.g. (not supported at the moment) $func_name(arg1, arg2, arg3, $func_name2(arg1, arg2, arg3))

	// check if it is function
	if len(key) < 3 {
		return "", errors.New("key is not a function")
	}

	if key[0] != '$' {
		return "", errors.New("key is not a function")
	}

	if key[len(key)-1] != ')' {
		return "", errors.New("key is not a function")
	}

	// check if it has '('
	pos := strings.Index(key, "(")
	if pos == -1 {
		return "", errors.New("key is not a function")
	}

	return key[1:pos], nil
}

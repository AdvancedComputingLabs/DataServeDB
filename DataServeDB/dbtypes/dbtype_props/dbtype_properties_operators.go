package dbtype_props

import "errors"

type NegateOperatorI interface {
	Negate()
}

// following function is used in parsing I don't think it needs to be public.

func GetNegatable(v interface{}) (NegateOperatorI, error) {

	if n, ok := v.(NegateOperatorI); ok {
		return n, nil
	}

	return nil, errors.New("field property '' of '' does not support negation")
}

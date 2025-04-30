package main

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var bucketColourMatchError = &tracer.Error{
	Kind: "bucketColourMatchError",
	Desc: "The number of buckets must match the number of colours.",
}

var stringToNumberError = &tracer.Error{
	Kind: "stringToNumberError",
	Desc: "The given string could not be converted to a number.",
}

func isStringToNumber(err error) bool {
	return errors.Is(err, stringToNumberError)
}

package main

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var outOfRangeError = &tracer.Error{
	Kind: "outOfRangeError",
	Desc: "The given number was not found to be between 0% and 100%.",
}

func isOutOfRange(err error) bool {
	return errors.Is(err, outOfRangeError)
}

var stringToNumberError = &tracer.Error{
	Kind: "stringToNumberError",
	Desc: "The given string could not be converted to a number.",
}

func isStringToNumber(err error) bool {
	return errors.Is(err, stringToNumberError)
}

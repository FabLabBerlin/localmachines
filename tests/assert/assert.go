package assert

import (
	"strconv"
)

func NoErrors(errs ...error) {
	msg := ""
	for i, err := range errs {
		if err != nil {
			msg += "[" + strconv.Itoa(i) + "] " + err.Error()
		}
	}
	if msg != "" {
		panic("Errors: " + msg)
	}
}

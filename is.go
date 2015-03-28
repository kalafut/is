package is

import (
	"fmt"
	"log"
	"testing"
)

// Is provides methods that leverage the existing testing capabilities found
// in the Go test framework. The methods provided allow for a more natural,
// efficient and expressive approach to writing tests. The goal is to write
// fewer lines of code while improving communication of intent.
type Is struct {
	TB         testing.TB
	fail       func(format string, args ...interface{})
	failFormat string
	failArgs   []interface{}
}

// New creates a new instance of the Is object and stores a reference to the
// provided testing object.
func New(tb testing.TB) *Is {
	if tb == nil {
		log.Fatalln("You must provide a testing object.")
	}
	v := &Is{TB: tb}
	v.fail = func(format string, args ...interface{}) {
		fmt.Print(decorate(fmt.Sprintf(format, args...)))
		if len(v.failFormat) != 0 {
			fmt.Printf("\t"+v.failFormat+"\n", v.failArgs...)
		}
		tb.FailNow()
	}
	return v
}

// Msg defines a message to print in the event of a failure. This allows you
// to print out additional information about a failure if it happens.
func (is *Is) Msg(format string, args ...interface{}) *Is {
	return &Is{
		TB:         is.TB,
		fail:       is.fail,
		failFormat: format,
		failArgs:   args,
	}
}

// Equal performs a deep compare of the provided objects and fails if they are
// not equal.
//
// Equal does not respect type differences. If the types are different and
// comparable (eg int32 and int64), but the values are the same, the objects
// are considered equal.
func (is *Is) Equal(a interface{}, b interface{}) {
	result := isEqual(a, b)
	if !result {
		is.fail("expected objects '%s' and '%s' to be equal, but got: %v and %v",
			objectTypeName(a),
			objectTypeName(b), a, b)
	}
}

// NotEqual performs a deep compare of the provided objects and fails if they are
// equal.
//
// NotEqual does not respect type differences. If the types are different and
// comparable (eg int32 and int64), but the values are different, the objects
// are considered not equal.
func (is *Is) NotEqual(a interface{}, b interface{}) {
	result := isEqual(a, b)
	if result {
		is.fail("expected objects '%s' and '%s' not to be equal",
			objectTypeName(a),
			objectTypeName(b))
	}
}

// Err checks the provided error object to determine if an error is present.
func (is *Is) Err(e error) {
	result := isNil(e)
	if result {
		is.fail("expected error")
	}
}

// NotErr checks the provided error object to determine if an error is not
// present.
func (is *Is) NotErr(e error) {
	result := isNil(e)
	if !result {
		is.fail("expected no error, but got: %v", e)
	}
}

// Nil checks the provided object to determine if it is nil.
func (is *Is) Nil(o interface{}) {
	result := isNil(o)
	if !result {
		is.fail("expected object '%s' to be nil, but got: %v", objectTypeName(o), o)
	}
}

// NotNil checks the provided object to determine if it is not nil.
func (is *Is) NotNil(o interface{}) {
	result := isNil(o)
	if result {
		is.fail("expected object '%s' not to be nil", objectTypeName(o))
	}
}

// True checks the provided boolean to determine if it is true.
func (is *Is) True(b bool) {
	result := b == true
	if !result {
		is.fail("expected boolean to be true")
	}
}

// False checks the provided boolean to determine if is false.
func (is *Is) False(b bool) {
	result := b == false
	if !result {
		is.fail("expected boolean to be false")
	}
}

// Zero checks the provided object to determine if it is the zero value
// for the type of that object. The zero value is the same as what the object
// would contain when initialized but not assigned.
//
// This method, for example, would be used to determine if a string is empty,
// an array is empty or a map is empty. It could also be used to determine if
// a number is 0.
//
// In cases such as slice, map, array and chan, a nil value is treated the
// same as an object with len == 0
func (is *Is) Zero(o interface{}) {
	result := isZero(o)
	if !result {
		is.fail("expected object '%s' to be zero value, but it was: %v", objectTypeName(o), o)
	}
}

// NotZero checks the provided object to determine if it is not the zero
// value for the type of that object. The zero value is the same as what the
// object would contain when initialized but not assigned.
//
// This method, for example, would be used to determine if a string is not
// empty, an array is not empty or a map is not empty. It could also be used
// to determine if a number is not 0.
//
// In cases such as slice, map, array and chan, a nil value is treated the
// same as an object with len == 0
func (is *Is) NotZero(o interface{}) {
	result := isZero(o)
	if result {
		is.fail("expected object '%s' not to be zero value", objectTypeName(o))
	}
}
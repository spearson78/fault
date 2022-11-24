package floc

import (
	"errors"
	"fmt"
	"runtime"
)

var Disable = false

type WithLocation struct {
	wrapped  error
	Location string
}

func (e *WithLocation) Error() string  { return fmt.Sprintf("<location> : %v", e.wrapped) }
func (e *WithLocation) Cause() error   { return e.wrapped }
func (e *WithLocation) Unwrap() error  { return e.wrapped }
func (e *WithLocation) String() string { return e.Error() }

func getLocation(depth int) string {
	if Disable {
		return ""
	}

	_, file, line, ok := runtime.Caller(depth)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	} else {
		return ""
	}

}

func WrapDepth(err error, depth int) error {
	if err == nil {
		return nil
	}

	loc := getLocation(depth + 2)

	return wrap(err, loc)
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	loc := getLocation(2)

	return wrap(err, loc)
}

func wrap(err error, location string) error {
	if err == nil {
		return nil
	}

	if location == "" {
		return err
	}

	return &WithLocation{
		wrapped:  err,
		Location: location,
	}
}

func WithDepth(depth int) func(error) error {

	loc := getLocation(depth + 2)

	return func(err error) error {
		return wrap(err, loc)
	}
}

func With() func(error) error {

	loc := getLocation(2)

	return func(err error) error {
		return wrap(err, loc)
	}
}

func Get(err error) ([]string, bool) {
	if err == nil {
		return nil, false
	}

	locations := make([]string, 0, 10)

	for err != nil {
		if withLocation, isLocation := err.(*WithLocation); isLocation {
			locations = append(locations, "")
			copy(locations[1:], locations)
			locations[0] = withLocation.Location
		}
		err = errors.Unwrap(err)
	}

	return locations, len(locations) > 0
}

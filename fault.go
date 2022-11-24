// Package fault provides an extensible yet ergonomic mechanism for wrapping
// errors. It implements this as a kind of middleware style pattern by providing
// a simple option-style interface that can be passed to a call to `fault.Wrap`.
//
// See the GitHub repository for full documentation and examples.
package fault

import (
	"fmt"
	"strings"

	"github.com/Southclaws/fault/floc"
)

// Wrapper describes a kind of middleware that packages can satisfy in order to
// decorate errors with additional, domain-specific context.
type Wrapper func(err error) error

var AutoGenerateLocations = true

// Wrap wraps an error with all of the wrappers provided.
func Wrap(err error, w ...Wrapper) error {
	if err == nil {
		return nil
	}

	generatedLocation := false

	for _, fn := range w {
		newErr := fn(err)
		if newErr != nil {
			err = newErr
		}

		if _, isLocation := newErr.(*floc.WithLocation); isLocation {
			generatedLocation = true
		}
	}

	if AutoGenerateLocations && !generatedLocation {
		err = floc.WrapDepth(err, 1)
	}

	return &container{
		cause: err,
	}
}

type container struct {
	cause error
}

// Error behaves like most error wrapping libraries, it gives you all the error
// messages conjoined with ": ". This is useful only for internal error reports,
// never show this to an end-user or include it in responses as it may reveal
// internal technical information about your application stack.
func (f *container) Error() string {
	errs := []string{}
	chain := Flatten(f)

	// reverse iterate since the chain is in caller order
	for i := len(chain.Errors) - 1; i >= 0; i-- {
		message := chain.Errors[i].Message
		if message != "" {
			errs = append(errs, chain.Errors[i].Message)
		}
	}

	return strings.Join(errs, ": ")
}

func (f *container) Unwrap() error { return f.cause }

func (f *container) Format(s fmt.State, verb rune) {
	u := Flatten(f)

	for _, v := range u.Errors {
		if v.Message != "" {
			fmt.Fprintf(s, "%s\n", v.Message)
		}
		if v.Location != "" {
			fmt.Fprintf(s, "\t%s\n", v.Location)
		}
	}
}

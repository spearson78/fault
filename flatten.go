package fault

import (
	"errors"

	"github.com/Southclaws/fault/floc"
)

// Chain represents an unwound error chain. It contains a reference to the root
// error (the error at the start of the chain, which does not provide any way
// to `errors.Unwrap` it further) as well as a slice of "Step" objects which are
// points where the error was wrapped.
type Chain struct {
	Root   error
	Errors []Step
}

// Step represents a location where an error was wrapped by Fault. The location
// is present if the error being wrapped contained stack information and the
// message is present if the underlying error provided a message. Note that not
// all errors provide errors or locations. If both are missing, it's omitted.
type Step struct {
	Location string
	Message  string
}

// Flatten attempts to derive more useful structured information from an error
// chain. If the input is a fault error, the output will contain an easy to use
// error chain list with location information and individual error messages.
func Flatten(err error) *Chain {
	if err == nil {
		return nil
	}

	if containerErr, firstErrIsContainer := err.(*container); firstErrIsContainer {
		err = containerErr.Unwrap()
	}

	var f Chain
	var s Step
	for err != nil {
		// Loop exits on the last error leaving root ass the last
		// error in the chain. This error is the root or "external" error.
		f.Root = err

		switch typedErr := err.(type) {
		case *container:
			//Add the current step to the chain and start a new step for this container
			f.Errors = append([]Step{s}, f.Errors...)
			s = Step{}
		case *floc.WithLocation:
			//Extract the first location from this step
			if s.Location == "" {
				s.Location = typedErr.Location
			}
		default:
			//Extract the first message from this step.
			if s.Message == "" {
				s.Message = typedErr.Error()
			}
		}

		err = errors.Unwrap(err)
	}

	//Add the last step to the chain
	f.Errors = append([]Step{s}, f.Errors...)

	return &f
}

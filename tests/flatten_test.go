package tests

import (
	"testing"

	"github.com/Southclaws/fault"
	"github.com/stretchr/testify/assert"
)

func TestFlattenStdlibSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(1)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: stdlib sentinel error", full)
	a.Equal("stdlib sentinel error", root)
	a.Len(chain.Errors, 3)

	e0 := chain.Errors[0]
	a.Equal("stdlib sentinel error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")

	e2 := chain.Errors[2]
	a.Equal("", e2.Message)
	a.Contains(e2.Location, "test_callers.go:11")
}

func TestFlattenFaultSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(2)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: fault sentinel error", full)
	a.Equal("fault sentinel error", root)
	a.Len(chain.Errors, 3)

	e0 := chain.Errors[0]
	a.Equal("fault sentinel error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")

	e2 := chain.Errors[2]
	a.Equal("", e2.Message)
	a.Contains(e2.Location, "test_callers.go:11")
}

func TestFlattenStdlibInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(3)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: stdlib root cause error", full)
	a.Equal("stdlib root cause error", root)
	a.Len(chain.Errors, 3)

	e0 := chain.Errors[0]
	a.Equal("stdlib root cause error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")

	e2 := chain.Errors[2]
	a.Equal("", e2.Message)
	a.Contains(e2.Location, "test_callers.go:11")
}

func TestFlattenFaultInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(4)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: fault root cause error", full)
	a.Equal("fault root cause error", root)
	a.Len(chain.Errors, 3)

	e0 := chain.Errors[0]
	a.Equal("fault root cause error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")

	e2 := chain.Errors[2]
	a.Equal("", e2.Message)
	a.Contains(e2.Location, "test_callers.go:11")
}

func TestFlattenStdlibErrorfWrappedError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(5)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: errorf wrapped: stdlib sentinel error", full)
	a.Equal("stdlib sentinel error", root)
	a.Len(chain.Errors, 3)

	e0 := chain.Errors[0]
	a.Equal("errorf wrapped: stdlib sentinel error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")

	e2 := chain.Errors[2]
	a.Equal("", e2.Message)
	a.Contains(e2.Location, "test_callers.go:11")
}

func TestFlattenStdlibErrorfWrappedExternalError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(6)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: errorf wrapped external: external error wrapped with errorf: stdlib external error", full)
	a.ErrorContains(err, "external error wrapped with errorf: stdlib external error")
	a.Equal("stdlib external error", root)
	a.Len(chain.Errors, 3)

	e0 := chain.Errors[0]
	a.Equal("errorf wrapped external: external error wrapped with errorf: stdlib external error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")

	e2 := chain.Errors[2]
	a.Equal("", e2.Message)
	a.Contains(e2.Location, "test_callers.go:11")
}

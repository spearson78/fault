package fault

import "fmt"

// New creates a new basic fault error.
func New(message string) error {
	return &fundamental{
		msg: message,
	}
}

// Newf includes formatting specifiers.
func Newf(message string, va ...any) error {
	return &fundamental{
		msg: fmt.Sprintf(message, va...),
	}
}

type fundamental struct {
	msg string
}

func (f *fundamental) Error() string {
	return f.msg
}

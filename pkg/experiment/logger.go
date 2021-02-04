package experiment

import "fmt"

// Logger represents the ability to log messages and errors.
// https://github.com/go-logr/logr example
type Logger interface {
	Enabled() bool
	Info(msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
}

// verify defaultLogger implements Logger interface
var _ Logger = (*defaultLogger)(nil)

type defaultLogger struct {
}

func (d *defaultLogger) Enabled() bool {
	return true
}

func (d *defaultLogger) Info(msg string, keysAndValues ...interface{}) {
	fmt.Println(fmt.Sprintf(msg, keysAndValues...))
}

func (d *defaultLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "error", err)
	d.Info(msg, keysAndValues...)
}


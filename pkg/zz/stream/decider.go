package stream

type Directive int

type Decider func(cause error) Directive

func StoppingDecider(cause error) Directive {
	return Stop
}

func ResumingDecider(cause error) Directive {
	return Resume
}

const (
	// The stream will be completed with failure if application code for processing an element returns an error.
	Stop Directive = iota
	// The element is dropped and the stream continues if application code for processing an element returns an error.
	Resume

	Reset
)

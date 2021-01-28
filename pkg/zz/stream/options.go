package stream

type StageOptions struct {
	Name             string
	OutputBufferSize int
	Parallelism      int
	Logger           Logger
	Decider Decider
}

func (s *StageOptions) Apply(options ...StageOption) *StageOptions {
	for _, apply := range options {
		apply(s)
	}
	return s
}

func (s *StageOptions) Copy() *StageOptions {
	return & StageOptions{
		Name:             s.Name,
		OutputBufferSize: s.OutputBufferSize,
		Parallelism:      s.Parallelism,
		Logger:           s.Logger,
	}
}

type StageOption func(*StageOptions)

var DefaultStageOptions = &StageOptions{
	Name:             "",
	OutputBufferSize: 0,
	Parallelism:      1,
	Logger:           &defaultLogger{},
	Decider: StoppingDecider,
}

func Name(value string) StageOption {
	return func(state *StageOptions) {
		state.Name = value
	}
}

func OutputBuffer(value int) StageOption {
	return func(state *StageOptions) {
		state.OutputBufferSize = value
	}
}

func Parallelism(value int) StageOption {
	return func(state *StageOptions) {
		state.Parallelism = value
	}
}

func ErrorStrategy(strategy Decider) StageOption {
	return func(state *StageOptions) {
		state.Decider = strategy
	}
}

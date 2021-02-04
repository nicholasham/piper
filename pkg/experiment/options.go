package experiment

type StageAttributes struct {
	Name             string
	OutputBufferSize int
	Parallelism      int
	Logger           Logger
	Decider          Decider
}

func (s *StageAttributes) Apply(options ...StageOption) *StageAttributes {
	for _, apply := range options {
		apply(s)
	}
	return s
}

func (s *StageAttributes) Copy() *StageAttributes {
	return &StageAttributes{
		Name:             s.Name,
		OutputBufferSize: s.OutputBufferSize,
		Parallelism:      s.Parallelism,
		Logger:           s.Logger,
	}
}

type StageOption func(*StageAttributes)

var DefaultStageAttributes = &StageAttributes{
	Name:             "",
	OutputBufferSize: 0,
	Parallelism:      1,
	Logger:           &defaultLogger{},
	Decider:          StoppingDecider,
}

func Name(value string) StageOption {
	return func(state *StageAttributes) {
		state.Name = value
	}
}

func OutputBuffer(value int) StageOption {
	return func(state *StageAttributes) {
		state.OutputBufferSize = value
	}
}

func Parallelism(value int) StageOption {
	return func(state *StageAttributes) {
		state.Parallelism = value
	}
}

func ErrorStrategy(strategy Decider) StageOption {
	return func(state *StageAttributes) {
		state.Decider = strategy
	}
}

package stream

type StageState struct {
	Name             string
	OutputBufferSize int
	Parallelism      int
	Logger           Logger
}

type StageOption func(*StageState)

var DefaultOptions = &StageState{
	Name:             "",
	OutputBufferSize: 0,
	Parallelism:      1,
	Logger:           &defaultLogger{},
}

func Name(value string) StageOption {
	return func(state *StageState) {
		state.Name = value
	}
}

func OutputBuffer(value int) StageOption {
	return func(state *StageState) {
		state.OutputBufferSize = value
	}
}

func Parallelism(value int) StageOption {
	return func(state *StageState) {
		state.Parallelism = value
	}
}

func NewStageState(name string, options ...StageOption) *StageState {
	state := DefaultOptions
	for _, apply := range options {
		apply(state)
	}
	return state
}

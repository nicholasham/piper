package piper

type OptionState struct {
	ChannelBuffer int
	Parallelism   int
}
type Option func(state *OptionState)

func Buffer(capacity int) Option {
	return func(state *OptionState) {
		state.ChannelBuffer = capacity
	}
}

func Parallelism(value int) Option {
	return func(state *OptionState) {
		state.Parallelism = value
	}
}

func getOptionState(options ...Option) *OptionState {
	state := &OptionState{}
	for _, option := range options {
		option(state)
	}
	return state
}

func CreateChannel(options []Option) chan Element {
	state := getOptionState(options...)
	if state.ChannelBuffer > 0 {
		return make(chan Element, state.ChannelBuffer)
	}
	return make(chan Element)
}

func GetParallelism(supported bool, options []Option) int {
	if supported {
		return getOptionState(options...).Parallelism
	}
	return 1
}

func Options(options ...Option) []Option {
	return options
}

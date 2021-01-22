package attribute

type StageAttributes struct {
	Name             string
	OutputBufferSize int
	Parallelism      int
	Logger           Logger
}

type StageAttribute func(stage *StageAttributes)

func Name(value string) StageAttribute {
	return func(state *StageAttributes) {
		state.Name = value
	}
}

func OutputBuffer(value int) StageAttribute {
	return func(state *StageAttributes) {
		state.OutputBufferSize = value
	}
}

func Parallelism(value int) StageAttribute {
	return func(state *StageAttributes) {
		state.Parallelism = value
	}
}

func Default(name string, attributes ...StageAttribute) *StageAttributes {

	stageAttributes := &StageAttributes{
		Name:             name,
		OutputBufferSize: 0,
		Parallelism:      1,
		Logger:           &defaultLogger{},
	}
	for _, attribute := range attributes {
		attribute(stageAttributes)
	}
	return stageAttributes
}

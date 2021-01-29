package experiment

import "github.com/nicholasham/piper/pkg/zz/stream"

type StageActions interface {
	CompleteStage()
	Push(element stream.Element)
	FailStage(cause error)
}

var _ StageActions = (*stageActions)(nil)

type stageActions struct {
	onCompleteStage func()
	onPush          func(element stream.Element)
	onFailStage     func(cause error)
}

func (s *stageActions) CompleteStage() {
	s.onCompleteStage()
}

func (s *stageActions) Push(element stream.Element) {
	s.onPush(element)
}

func (s *stageActions) FailStage(cause error) {
	s.onFailStage(cause)
}

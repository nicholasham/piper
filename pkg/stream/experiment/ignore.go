package experiment

import "github.com/nicholasham/piper/pkg/stream"

var _ InHandler = (*IgnoreSinkLogic)(nil)


type IgnoreSinkLogic struct {

}

func (i *IgnoreSinkLogic) OnPush(element stream.Element) {
}

func (i *IgnoreSinkLogic) OnUpstreamFinish(actions StageActions) {
	actions.CompleteStage()
}

func (i *IgnoreSinkLogic) OnUpstreamFailure(cause error, actions StageActions) {
	actions.FailStage(cause)
}

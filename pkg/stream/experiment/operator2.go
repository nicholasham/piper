package experiment

import (
	"github.com/nicholasham/piper/pkg/stream"
)

type InOutLogic interface {
	InHandler
	OutHandler
}


type InHandler interface {
	OnPush(element stream.Element)
	OnUpstreamFinish(actions StageActions)
	OnUpstreamFailure(cause error, actions StageActions)
}


type OutHandler interface {
	OnPull(actions StageActions)
	OnDownstreamFinish(actions StageActions)
}


package experiment

import (
	"github.com/nicholasham/piper/pkg/zz/stream"
)

type InOutLogic interface {
	InHandler
	OutHandler
}


type InHandler interface {
	OnPush(element stream.Element, actions StageActions)
	OnUpstreamFinish(actions StageActions)
	OnUpstreamFailure(cause error, actions StageActions)
}


type OutHandler interface {
	OnPull( actions StageActions)
	OnDownstreamFinish(actions StageActions)
}


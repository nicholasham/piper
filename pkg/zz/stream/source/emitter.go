package source

import (
	"context"
)

type Emitter interface {
	Run(ctx context.Context, actions EmitterActions)
}

type EmitterActions interface {
}

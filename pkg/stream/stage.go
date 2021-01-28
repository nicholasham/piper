package stream

import "context"

type stage struct {
	outlets []*Outlet
	inlets []*Inlet
	strategy Strategy
}

type StrategyFactory func() Strategy

type Strategy func(ctx context.Context, inlets []*Inlet, outlet[] *Outlet)

func (s *stage) Run(ctx context.Context) {
	go s.strategy(ctx, s.inlets, s.outlets)
}


type InOutStrategy func(ctx context.Context, inlet *Inlet, outlet *Outlet)



func nick(inOut InOutStrategy) Strategy {
	return func(ctx context.Context, inlets []*Inlet, outlet []*Outlet) {
		inOut(ctx, inlets[0], outlet[0])
	}
}

func Test(ctx context.Context, inlet *Inlet, outlet *Outlet)  {
	for element := range inlet.in {
		outlet.Send(element)
	}
}


package runners

import (
	"api_crypto1.0/internal/usecases"
	"context"
	"sync"
)

type Runner struct {
	Ucase *usecases.Usecases
	Ctx   context.Context
	Wg    *sync.WaitGroup
}

func NewRunner(usecase *usecases.Usecases, ctx context.Context, wg *sync.WaitGroup) *Runner {
	return &Runner{
		Ucase: usecase,
		Ctx:   ctx,
		Wg:    wg,
	}
}

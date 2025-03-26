package run

import "context"

type RunUseCaseInterface interface {
	Run(ctx context.Context, input RunInputDTO) (RunOutputDTO, error)
}

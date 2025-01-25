package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

func FillEnvStruct(c any) error {
	ctx := context.Background()
	return envconfig.Process(ctx, c)
}

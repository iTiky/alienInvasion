package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/itiky/alienInvasion/cmd/alieninvasion"
	"github.com/itiky/alienInvasion/pkg/logging"
)

func main() {
	rand.Seed(time.Now().Unix())

	if err := alieninvasion.NewRootCmd().Execute(); err != nil {
		_, logger := logging.GetCtxLogger(context.Background())
		logger.Fatal().Err(err).Msg("Application failed")
	}
}

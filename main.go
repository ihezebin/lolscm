package main

import (
	"context"
	"lolscm/cmd"

	"github.com/ihezebin/oneness/logger"
)

func main() {
	ctx := context.Background()
	if err := cmd.Run(ctx); err != nil {
		logger.Fatalf(ctx, "lolscm run error: %v", err)
	}
}

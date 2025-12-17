package main

import (
	"github.com/scmbr/oms/common/logger"
	"github.com/scmbr/oms/order-service/app"
)

func main() {
	if err := app.Run(); err != nil {
		logger.Error("application terminated with error", err)
	}
}

package healthcheck

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestChecker(t *testing.T) {
	healthFn := func() error {
		fmt.Println("health")
		return nil
	}

	unHealthFn := func() error {
		fmt.Println("unHealth")
		return nil
	}

	HealthcheckFile = "/tmp/healthcheck.html"
	logger, _ := zap.NewProduction()
	c := NewChecker(time.Second, healthFn, unHealthFn, logger)
	go c.CheckForever()

	stop := make(chan struct{})
	<-stop
}

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

	logger, _ := zap.NewProduction()
	c := NewChecker(time.Second, NewFileHealth("/tmp/healthcheck.html"), healthFn, unHealthFn, logger)
	go c.CheckForever()

	stop := make(chan struct{})
	<-stop
}

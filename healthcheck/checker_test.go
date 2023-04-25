package healthcheck

import (
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func TestChecker(t *testing.T) {
	healthFn := func() error {
		log.Println("health")
		return nil
	}

	unHealthFn := func() error {
		log.Println("unHealth")
		return nil
	}

	logger, _ := zap.NewProduction()
	c := NewChecker(time.Second, NewFileHealth("/tmp/healthcheck.html"), healthFn, unHealthFn, logger)
	go c.CheckForever()

	stop := make(chan struct{})
	<-stop
}

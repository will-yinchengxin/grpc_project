package code

import (
	"errors"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
	"log"
	"math/rand"
	"testing"
	"time"
)

// https://github.com/alibaba/sentinel-golang/blob/master/example/circuitbreaker/error_count/circuit_breaker_error_count_example.go

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func TestCircuit_ErrCount(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:       "abc",
			Strategy:       circuitbreaker.ErrorCount,
			Threshold:      50,
			StatIntervalMs: 5000,

			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,

			StatSlidingWindowBucketCount: 10,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	logging.Info("[CircuitBreaker ErrorCount] Sentinel Go circuit breaking demo is running. You may see the pass/block metric in the metric log.")

	go func() {
		for {
			e, b := sentinel.Entry("abc")
			if b != nil {
				// g1 blocked
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				if rand.Uint64()%20 > 9 {
					// Record current invocation as error.
					sentinel.TraceError(e, errors.New("biz error"))
				}
				// g1 passed
				time.Sleep(time.Duration(rand.Uint64()%80+10) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	go func() {
		for {
			e, b := sentinel.Entry("abc")
			if b != nil {
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(rand.Uint64()%80) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	<-ch
}

func TestCircuit_ErrRatio(t *testing.T) {

}

func TestCircuit_SlowRTRatio(t *testing.T) {

}

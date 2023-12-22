package code

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/util"
	"log"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

type Counter struct {
	pass  *int64
	block *int64
	total *int64
}

var routineCount = 30

func TestWarmUp(t *testing.T) {
	counter := Counter{pass: new(int64), block: new(int64), total: new(int64)}
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.WarmUp,
			ControlBehavior:        flow.Reject,
			Threshold:              20,
			StatIntervalInMs:       1000,
			WarmUpPeriodSec:        10,
			WarmUpColdFactor:       3,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
	go timerTask(&counter)
	ch := make(chan struct{})
	//warmUp task
	for i := 0; i < 3; i++ {
		go Task(&counter)
	}
	time.Sleep(3 * time.Second)
	////sentinel task
	//for i := 0; i < routineCount; i++ {
	//	go Task(&counter)
	//}
	<-ch
}

func Task(counter *Counter) {
	for {
		atomic.AddInt64(counter.total, 1)
		e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			atomic.AddInt64(counter.block, 1)
		} else {
			// Be sure the entry is exited finally.
			e.Exit()
			atomic.AddInt64(counter.pass, 1)
		}
		time.Sleep(time.Duration(rand.Uint64()%50) * time.Millisecond)
	}
}
func timerTask(counter *Counter) {
	fmt.Println("begin to statistic!!!")
	var (
		oldTotal, oldPass, oldBlock int64
	)
	for {
		time.Sleep(1 * time.Second)
		globalTotal := atomic.LoadInt64(counter.total)
		oneSecondTotal := globalTotal - oldTotal
		oldTotal = globalTotal

		globalPass := atomic.LoadInt64(counter.pass)
		oneSecondPass := globalPass - oldPass
		oldPass = globalPass

		globalBlock := atomic.LoadInt64(counter.block)
		oneSecondBlock := globalBlock - oldBlock
		oldBlock = globalBlock

		fmt.Println("time:", util.CurrentTimeMillis()/1000, "total:", oneSecondTotal, " pass:", oneSecondPass, " block:", oneSecondBlock)
	}
}

func TestFlowTwo(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			Threshold:              5,
			StatIntervalInMs:       1000,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					fmt.Println(util.CurrentTimeMillis(), "Reject")
					time.Sleep(time.Second * 2)
				} else {
					fmt.Println(util.CurrentTimeMillis(), "Passed")
					//time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					time.Sleep(time.Second * 2)
					e.Exit()
				}

			}
		}()
	}
	<-ch

	/*
		1702016951106 Passed
		1702016951106 Passed
		1702016951106 Passed
		1702016951106 Passed
		1702016951106 Passed
		1702016951106 Reject
		1702016951106 Reject
		1702016951106 Reject
		1702016951106 Reject
		1702016951106 Reject
	*/
}

func TestFlowOne(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			Threshold:              5,
			StatIntervalInMs:       1000,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					fmt.Println(util.CurrentTimeMillis(), "Reject")
					time.Sleep(time.Second * 2)
				} else {
					fmt.Println(util.CurrentTimeMillis(), "Passed")
					//time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					time.Sleep(time.Second * 2)
					e.Exit()
				}

			}
		}()
	}
	<-ch

	/*
		1702016951106 Passed
		1702016951106 Passed
		1702016951106 Passed
		1702016951106 Passed
		1702016951106 Passed
		1702016951106 Reject
		1702016951106 Reject
		1702016951106 Reject
		1702016951106 Reject
		1702016951106 Reject
	*/
}

func TestSentinel(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	// initialize sentinel rules
	if err = initRules(); err != nil {
		t.Fatal(err)
		return
	}

	ch := make(chan struct{})
	for i := 0; i < 5; i++ {
		go func() {
			for {
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					fmt.Println(util.CurrentTimeMillis(), "Reject")
					time.Sleep(time.Second * 2)
				} else {
					fmt.Println(util.CurrentTimeMillis(), "Passed")
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)

					e.Exit()
				}

			}
		}()
	}
	<-ch
}

func initRules() error {
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			Threshold:              10,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	if err != nil {
		return err
	}
	// 可以通过各个模块的 LoadRules(rules) 函数加载规则。
	//_, err = hotspot.LoadRules([]*hotspot.Rule{
	//	{
	//		Resource:        "some-test",
	//		Threshold:       10,
	//		ControlBehavior: hotspot.Reject,
	//	},
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}

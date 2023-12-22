# [Sentinel](https://sentinelguard.io/zh-cn/docs/introduction.html)
Sentinel 是面向分布式、多语言异构化服务架构的流量治理组件，主要以流量为切入点，
从流量路由、流量控制、流量整形、熔断降级、系统自适应过载保护、热点流量防护等多个维度来帮助开发者保障微服务的稳定性。

## [Sentinel-Go](https://sentinelguard.io/zh-cn/docs/golang/quick-start.html)
Sentinel 中所有的限流熔断机制都是基于资源生效的，不同资源的限流熔断规则互相隔离互不影响。

Demo 可见 `./code/sentinel_test.go`

### API 基本使用规则
#### 通用配置及初始化
必须成功调用 Sentinel 的初始化函数以后再调用埋点 API !!!
```go
package main

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
)

func initSentinel() {
	err := sentinel.Init(confPath)
	if err != nil {
		// 初始化 Sentinel 失败
	}
}
````
#### 埋点(定义资源)
Sentinel 的 Entry API 将业务逻辑封装起来，这一步称为“埋点”。

``Entry(resource string, opts ...Option) (*base.SentinelEntry, *base.BlockError)``

resource 代表埋点资源名，opts 代表埋点配置。这里需要注意的是，返回值参数列表的第一个和第二个参数是互斥的，也就是说，如果Entry执行pass，
那么Sentinel会返回(*base.SentinelEntry, nil)；如果Entry执行blocked，那么Sentinel会返回(nil, *base.BlockError)。

#### 规则配置

Sentinel 支持原始的硬编码方式加载规则，可以通过各个模块的 LoadRules(rules) 函数加载规则。以流控规则为例：

```go
_, err = flow.LoadRules([]*flow.Rule{
	{
		Resource:               "some-test",
		Threshold:              10,
		TokenCalculateStrategy: flow.Direct,
		ControlBehavior:        flow.Reject,
	},
})
if err != nil {
	// 加载规则失败，进行相关处理
}
````

## [流量控制 (Sentinel Go)](https://sentinelguard.io/zh-cn/docs/golang/flow-control.html)

```go
type Rule struct {
	ID string `json:"id,omitempty"`
	
	Resource               string                 `json:"resource"`
	TokenCalculateStrategy TokenCalculateStrategy `json:"tokenCalculateStrategy"`
	ControlBehavior        ControlBehavior        `json:"controlBehavior"`
	Threshold        float64          `json:"threshold"`
	RelationStrategy RelationStrategy `json:"relationStrategy"`
	RefResource      string           `json:"refResource"`
	MaxQueueingTimeMs uint32 `json:"maxQueueingTimeMs"`
	WarmUpPeriodSec   uint32 `json:"warmUpPeriodSec"`
	WarmUpColdFactor  uint32 `json:"warmUpColdFactor"`
	StatIntervalInMs uint32 `json:"statIntervalInMs"`
	
	// 以下参数可以暂且忽律
    LowMemUsageThreshold  int64 `json:"lowMemUsageThreshold"`
	HighMemUsageThreshold int64 `json:"highMemUsageThreshold"`
	MemLowWaterMarkBytes  int64 `json:"memLowWaterMarkBytes"`
	MemHighWaterMarkBytes int64 `json:"memHighWaterMarkBytes"`
}
````
一条流控规则主要由下面几个因素组成，我们可以组合这些元素来实现不同的限流效果：

- Resource：资源名，即规则的作用目标。
- StatIntervalInMs: 规则对应的流量控制器的 `独立统计结构的统计周期`。如果StatIntervalInMs是1000，也就是统计QPS。默认为 1s。
- Threshold: 表示 `流控阈值`；如果字段 StatIntervalInMs 是1000(也就是1秒)，那么Threshold就表示QPS，流量控制器也就会依据资源的QPS来做流控。
- TokenCalculateStrategy: 当前流量控制器的Token计算策略。`Direct`表示直接使用字段 Threshold(门槛) 作为阈值；`WarmUp`表示使用预热方式计算Token的阈值。
- ControlBehavior: 表示流量控制器的控制策略；`Reject`表示超过 阈值 直接拒绝，`Throttling`表示匀速排队。
- WarmUpPeriodSec: 预热的时间长度，该字段仅仅对 `WarmUp` 的 `TokenCalculateStrategy` 生效；
- WarmUpColdFactor: 预热的因子，默认是3，该值的设置会影响预热的速度，该字段仅仅对 `WarmUp` 的 `TokenCalculateStrategy` 生效；

- MaxQueueingTimeMs: 匀速排队的最大等待时间，该字段仅仅对 Throttling ControlBehavior生效；
- RefResource: 关联的resource；
- RelationStrategy: 调用关系限流策略，CurrentResource表示使用当前规则的resource做流控；AssociatedResource表示使用关联的resource做流控，
  关联的resource在字段 RefResource 定义；

Sentinel 的流量控制策略由规则中的 TokenCalculateStrategy 和 ControlBehavior 两个字段决定。
- TokenCalculateStrategy 表示流量控制器的Token计算方式
- ControlBehavior 表示流量控制器的控制行为

每 100ms 最多通过一个请求，多余的请求将会排队等待通过，若排队时队列长度大于 500ms 则直接拒绝：
```go
{
	Resource:          "some-test",
    TokenCalculateStrategy: flow.Direct,
	ControlBehavior:   flow.Throttling, // 流控效果为匀速排队
    Threshold:             10,          // 请求的间隔控制在 1000/10 = 100ms
	MaxQueueingTimeMs: 500,             // 最长排队等待时间
}

// Threshold 是 10，Sentinel 默认使用1s作为控制周期，表示1秒内10个请求匀速排队，所以排队时间就是 1000ms/10 = 100ms；
// MaxQueueingTimeMs 设为 0 时代表不允许排队，只控制请求时间间隔，多余的请求将会直接拒绝。
````
### 常见场景的规则配置-QPS
**基于对某个资源访问的QPS来做流控**
```go
{
  Resource:                "some-test",
  TokenCalculateStrategy:  flow.Direct,
  ControlBehavior:         flow.Reject,
  Threshold:               500,
  StatIntervalInMs:        1000,
}
// StatIntervalInMs必须是1000，表示统计周期是1s，那么Threshold所配置的值也就是QPS的阈值。
````
```go
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
````
**基于一定统计间隔时间来控制总的请求数**
比如StatIntervalInMs配置10000，Threshold配置10000，这种配置意思就是控制10s内最大请求数是10000。
```go
{
	Resource:                "some-test",
    TokenCalculateStrategy:  flow.Direct,
	ControlBehavior:         flow.Reject,
	Threshold:               10000,
    StatIntervalInMs:        10000,
}

// 比如StatIntervalInMs配置10000，Threshold配置10000，这种配置意思就是控制10s内最大请求数是10000
// 这种流控配置对于脉冲类型的流量抵抗力很弱，有极大潜在风险压垮系统。
````

**毫秒级别流控**
```go
{
	Resource:                "some-test",
    TokenCalculateStrategy:  flow.Direct,
	ControlBehavior:         flow.Reject,
	Threshold:               80,
    StatIntervalInMs:        100,
}

// 限制了100ms的阈值是80，实际QPS大概是800。
// 这种配置也是有缺点的，脉冲流量很大可能造成有损(会拒绝很多流量)。
````
### 常见场景的规则配置-WarmUp

```go
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
			Threshold:              10,
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
	for i := 0; i < 3; i++ {
		go Task(&counter)
	}
	time.Sleep(3 * time.Second)
	for i := 0; i < routineCount; i++ {
		go Task(&counter)
	}
	<-ch

	/*
		begin to statistic!!!
		1702020164 total: 126  pass: 6  block: 119
		1702020165 total: 133  pass: 3  block: 131
		1702020166 total: 159  pass: 3  block: 156
		1702020167 total: 1311  pass: 3  block: 1308
		1702020168 total: 1349  pass: 4  block: 1345
		1702020169 total: 1341  pass: 4  block: 1337
		1702020170 total: 1318  pass: 4  block: 1314
		1702020171 total: 1314  pass: 5  block: 1309
		1702020172 total: 1345  pass: 5  block: 1340
		1702020173 total: 1352  pass: 6  block: 1346
		1702020174 total: 1378  pass: 7  block: 1371
		1702020175 total: 1333  pass: 10  block: 1323
		1702020176 total: 1359  pass: 10  block: 1349
		1702020177 total: 1328  pass: 10  block: 1318
	*/
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
		fmt.Println(util.CurrentTimeMillis()/1000, "total:", oneSecondTotal, " pass:", oneSecondPass, " block:", oneSecondBlock)
	}
}
````
## [熔断降级](https://sentinelguard.io/zh-cn/docs/golang/circuit-breaking.html)
不同服务之间相互调用，组成复杂的调用链路。前面描述的问题在分布式链路调用中会产生放大的效果。整个复杂链路中的某一环如果不稳定，就可能会层层级联，
最终可能导致整个链路全部挂掉。因此我们需要对不稳定的 弱依赖服务调用 进行 熔断降级，暂时切断不稳定的服务调用，避免局部不稳定因素导致整个分布式系统的雪崩。
`熔断降级作为保护服务自身的手段，通常在客户端（调用端）进行配置。`


熔断器有三种状态：
- `Closed状态`：也是初始状态，该状态下，熔断器会保持闭合，对资源的访问直接通过熔断器的检查。
- `Open状态`：断开状态，熔断器处于开启状态，对资源的访问会被切断。
- `Half-Open状态`：半开状态，该状态下除了探测流量，其余对资源的访问也会被切断。探测流量指熔断器处于半开状态时，会周期性的允许一定数目的探测请求通过，
  如果探测请求能够正常的返回，代表探测成功，此时熔断器会重置状态到 Closed 状态，结束熔断；如果探测失败，则回滚到 Open 状态。

Sentinel 提供了监听器去监听熔断器状态机的三种状态的转换
```go
type StateChangeListener interface {
	OnTransformToClosed(prev State, rule Rule)
	OnTransformToOpen(prev State, rule Rule, snapshot interface{})
	OnTransformToHalfOpen(prev State, rule Rule)
}
````

我们衡量下游服务质量时候，场景的指标就是`RT(response time)`、`异常数量`以及`异常比例`等等

Sentinel 的熔断器支持三种熔断策略：`慢调用比例熔断`、`异常比例熔断` 以及 `异常数量熔断`。
- 慢调用比例策略 (SlowRequestRatio)：Sentinel 的熔断器不在`静默期`，并且慢调用的比例大于设置的阈值，则接下来的熔断周期内对资源的访问会自动地被熔断。
  该策略下需要设置允许的调用 RT 临界值（即最大的响应时间），对该资源访问的响应时间大于该阈值则统计为慢调用。
- 错误比例策略 (ErrorRatio)：Sentinel 的熔断器不在`静默期`，并且在统计周期内资源请求访问异常的比例大于设定的阈值，则接下来的熔断周期内对资源的访问会自动地被熔断。
- 错误计数策略 (ErrorCount)：Sentinel 的熔断器不在`静默期`，并且在统计周期内资源请求访问异常数大于设定的阈值，则接下来的熔断周期内对资源的访问会自动地被熔断。


```go
type Rule struct {
	Id string `json:"id,omitempty"`
	Resource string   `json:"resource"`
	Strategy Strategy `json:"strategy"`
	RetryTimeoutMs uint32 `json:"retryTimeoutMs"`
	MinRequestAmount uint64 `json:"minRequestAmount"`
	StatIntervalMs uint32 `json:"statIntervalMs"`
	MaxAllowedRtMs uint64 `json:"maxAllowedRtMs"`
	Threshold float64 `json:"threshold"`
}
````
- Resource: 熔断器规则生效的埋点资源的名称；
- Strategy: 熔断策略，目前支持SlowRequestRatio、ErrorRatio、ErrorCount三种；
- - 选择以慢调用比例 (SlowRequestRatio) 作为阈值，需要设置允许的最大响应时间（MaxAllowedRtMs），请求的响应时间大于该值则统计为慢调用。
   通过 Threshold 字段设置触发熔断的慢调用比例，取值范围为 [0.0, 1.0]。规则配置后，在单位统计时长内请求数目大于设置的最小请求数目，
   并且慢调用的比例大于阈值，则接下来的熔断时长内请求会自动被熔断。经过熔断时长后熔断器会进入探测恢复状态，若接下来的一个请求响应时间小于设置的最大 RT 则结束熔断， 
   若大于设置的最大 RT 则会再次被熔断。
- - 选择以错误比例 (ErrorRatio) 作为阈值，需要设置触发熔断的异常比例（Threshold），取值范围为 [0.0, 1.0]。
    规则配置后，在单位统计时长内请求数目大于设置的最小请求数目，并且异常的比例大于阈值，则接下来的熔断时长内请求会自动被熔断。
    经过熔断时长后熔断器会进入探测恢复状态，若接下来的一个请求没有错误则结束熔断，否则会再次被熔断。代码中可以通过 api.TraceError(entry, err) 函数来记录 error。
  
- RetryTimeoutMs: 即熔断触发后持续的时间（单位为 ms）。资源进入熔断状态后，在配置的熔断时长内，请求都会快速失败。熔断结束后进入探测恢复模式（HALF-OPEN）。
- MinRequestAmount: 静默数量，如果当前统计周期内对资源的访问数量小于静默数量，那么熔断器就处于静默期。换言之，也就是触发熔断的最小请求数目
  若当前统计周期内的请求数小于此值，即使达到熔断条件规则也不会触发。
- StatIntervalMs: 统计的时间窗口长度（单位为 ms）。
- MaxAllowedRtMs: 仅对慢调用熔断策略生效，MaxAllowedRtMs 是判断请求是否是慢调用的临界值，也就是如果请求的response time小于或等于MaxAllowedRtMs，那么就不是慢调用；
  如果response time大于MaxAllowedRtMs，那么当前请求就属于慢调用。
- Threshold: 对于慢调用熔断策略, Threshold表示是慢调用比例的阈值(小数表示，比如0.1表示10%)，也就是如果当前资源的慢调用比例如果高于Threshold，
  那么熔断器就会断开；否则保持闭合状态。 对于错误比例策略，Threshold表示的是错误比例的阈值(小数表示，比如0.1表示10%)。对于错误数策略，Threshold是错误计数的阈值。

```go
// 慢调用比例规则
rule1 := &Rule{
    Resource:         "abc",
    Strategy:         SlowRequestRatio,
    MaxAllowedRtMs:   20,    // 允许的最大响应时间
    Threshold:        0.1,   // 通过 Threshold 字段设置触发熔断的慢调用比例，取值范围为 [0.0, 1.0]
    StatIntervalMs:   10000, // 统计的时间窗口长度
	RetryTimeoutMs:   5000,  // 即熔断触发后持续的时间
	MinRequestAmount: 10,    // 静默数量
},
// 错误比例规则
rule1 := &Rule{
    Resource:         "abc",
    Strategy:         ErrorRatio,
	Threshold:        0.1,
    StatIntervalMs:   10000,
	RetryTimeoutMs:   5000,
	MinRequestAmount: 10,
},
// 错误计数规则
rule1 := &Rule{
    Resource:         "abc",
    Strategy:         ErrorCount,
	Threshold:        100,
    StatIntervalMs:   10000,
	RetryTimeoutMs:   5000,
	MinRequestAmount: 10,
},
````
```go
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
````


.....
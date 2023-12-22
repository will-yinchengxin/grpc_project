package rocketMQ

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	groupName = "Test_MQ"
	topic     = "Test_MQ_Topic"
	//mqAddr    = "172.16.27.95:9876"
	mqAddr = "127.0.0.1:9876"
)

// 出现以下问题, 检查以下网络, 当前宿主机是否能够 ping 通 172.20.0.4
// time="2023-12-11T15:04:39+08:00" level=warning msg="send heart beat to broker error" underlayError="dial tcp 172.20.0.4:10911: i/o timeout"
// 2023/12/11 15:04:44 Send Msg Err: dial tcp 172.20.0.4:10911: i/o timeout

func GetMqAddr() string {
	//mqAddr := fmt.Sprintf("%s:%d", internal.AppConf.RocketMQConfig.Host,
	//	internal.AppConf.RocketMQConfig.Port)
	return mqAddr
}

func ProduceMsg(mqAddr string, topic string) {
	// 普通消息生产者
	p, err := rocketmq.NewProducer(
		producer.WithGroupName(groupName),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{mqAddr})),
		producer.WithRetry(2),
	)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		log.Fatal("producer err: " + err.Error())
		os.Exit(1)
	}
	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("MQ-Test-----" + strconv.Itoa(i)),
		}
		msg.WithDelayTimeLevel(3)
		r, err := p.SendSync(context.Background(), msg)
		if err != nil {
			log.Fatal("Send Msg Err: " + err.Error())
		} else {
			log.Println("Send Msg Success: " + r.String() + "-" + r.MsgID)
		}
	}
	err = p.Shutdown()
	if err != nil {
		log.Fatal("producer shutdown" + err.Error())
		os.Exit(1)
	}
}

func ConsumeMsg(mqAddr string, topic string) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{mqAddr})),
	)
	if err != nil {
		panic(err)
	}
	err = c.Subscribe(topic, consumer.MessageSelector{},
		func(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgList {
				fmt.Printf("订阅消息，消费%v \n", msgList[i])
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		log.Fatal("consumer msg err: " + err.Error())
	}
	err = c.Start()
	if err != nil {
		log.Fatal("start consume err: " + err.Error())
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		log.Fatal("shutdown consume err: " + err.Error())
	}

}

func TestMQ(t *testing.T) {
	ProduceMsg(GetMqAddr(), topic)
	//ConsumeMsg(GetMqAddr(), topic)
	t.Log("test ok")
}

package rocketMQ

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"testing"
	"time"
)

func TestTrans(t *testing.T) {
	mqAddr := "127.0.0.1:9876"
	p, err := rocketmq.NewTransactionProducer(
		Trans{},
		producer.WithNameServer([]string{mqAddr}),
	)
	if err != nil {
		panic(err) // 生产环境禁用panic
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	res, err := p.SendMessageInTransaction(context.Background(),
		primitive.NewMessage("TransactionTopic", []byte("TransactionTopic Test")))
	fmt.Println(res.Status)
	if err != nil {
		panic(err)
	}
	fmt.Printf("发送成功")
	time.Sleep(time.Second * 3600)
	err = p.Shutdown()
	if err != nil {
		panic(err)
	}
}

type Trans struct {
}

func (hl Trans) ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

func (hl Trans) CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

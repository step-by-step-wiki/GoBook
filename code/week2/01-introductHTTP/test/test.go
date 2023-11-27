package main

import (
	"context"
	"crypto/tls"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"strconv"
	"sync"
	"time"
)

type ProducerMessage struct {
	Seq       int64
	Value     []byte
	Metadata  interface{}
	Timestamp time.Time
}
type AsyncProducer interface {
	Close() error
	Input() chan<- *ProducerMessage
	Successes() <-chan *ProducerMessage
	Errors() <-chan *ProducerMessage
}
type ProducerConfig struct {
	Endpoints    []string
	DialOptions  []grpc.DialOption
	BufferLength int
	TlsOptions   *tls.Config
}
type asyncProducer struct {
	client      *clientv3.Client
	wg          *sync.WaitGroup
	inputChan   chan *ProducerMessage
	bufChan     chan *ProducerMessage
	successChan chan *ProducerMessage
	errChan     chan *ProducerMessage
	err         error
	config      ProducerConfig
}

func NewAsyncProducer(config ProducerConfig) (AsyncProducer, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialOptions: config.DialOptions,
		TLS:         config.TlsOptions,
	})
	if err != nil {
		return nil, err
	}
	ap := &asyncProducer{
		client:      client,
		wg:          new(sync.WaitGroup),
		inputChan:   make(chan *ProducerMessage),
		bufChan:     make(chan *ProducerMessage, config.BufferLength),
		successChan: make(chan *ProducerMessage),
		errChan:     make(chan *ProducerMessage),
		err:         nil,
		config:      config,
	}
	go ap.run()
	return ap, nil
}
func (ap *asyncProducer) Input() chan<- *ProducerMessage {
	//randNum := genRandInt()
	//
	//if randNum > 50 {
	//	return ap.successChan
	//} else {
	//	return ap.errChan
	//}

	return ap.inputChan
}

//func genRandInt() int {
//	return rand.Intn(101)
//}

func (ap *asyncProducer) Successes() <-chan *ProducerMessage {
	return ap.successChan
}
func (ap *asyncProducer) Errors() <-chan *ProducerMessage {
	return ap.errChan
}
func (ap *asyncProducer) Close() error {
	close(ap.inputChan)
	ap.wg.Wait()
	return ap.err
}

// TODO, please write this function to handle message
func (ap *asyncProducer) run() {}
func (ap *asyncProducer) handleMessage(msg *ProducerMessage) {
	defer ap.wg.Done()
	seq, err := putSequential(ap.client, "test", string(msg.Value))
	msg.Timestamp = time.Now()
	if err != nil {
		ap.handleMessageError(msg, err)
		return
	}
	msg.Seq = seq
	ap.successChan <- msg
}
func (ap *asyncProducer) handleMessageError(msg *ProducerMessage, err error) {
	ap.errChan <- msg
}

// putSequential will get the sequence number of the current queue, add 1 and generate a
// new KV pairs. Key will be the queues/{name}/pending/sequence instead of etcd timestamp.
// Value will be the Message body. This method is mainly used to make KV pairs more human friendly.
func putSequential(client *clientv3.Client, name, val string) (int64, error) {
	seqKey := "test-seq"
	resp, err := client.Get(context.TODO(), seqKey)
	if err != nil {
		return 0, err
	}
	var seqNum int64
	if len(resp.Kvs) != 0 {
		seqNum, err = strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
		if err != nil {
			return 0, err
		}
		seqNum++
	}
	penKey := "test-pen"
	cmp := clientv3.Compare(clientv3.ModRevision(seqKey), "<", resp.Header.Revision+1)
	reqIncSeq := clientv3.OpPut(seqKey, strconv.FormatInt(seqNum, 10))
	reqPutPen := clientv3.OpPut(penKey, val)
	txnResp, err := client.Txn(context.TODO()).If(cmp).Then(reqIncSeq, reqPutPen).Commit()
	if err != nil {
		return 0, err
	}
	if !txnResp.Succeeded {
		return putSequential(client, name, val)
	}
	return seqNum, nil
}

func main() {
	ap, err := NewAsyncProducer(ProducerConfig{
		Endpoints:    []string{"127.0.0.1:2379"},
		DialOptions:  []grpc.DialOption{grpc.WithInsecure()},
		BufferLength: 10,
	})
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case msg := <-ap.Successes():
				fmt.Println(string(msg.Value), "success")
			case msg := <-ap.Errors():
				fmt.Println(string(msg.Value), msg.Seq, "error")
			default:
				fmt.Println("no message")
			}
		}
	}()
	for i := 0; i < 30; i++ {
		ap.Input() <- &ProducerMessage{Value: []byte(strconv.Itoa(i))}
		<-time.After(time.Second)
	}
	err = ap.Close()
}

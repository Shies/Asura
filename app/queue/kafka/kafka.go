package kafka

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"Asura/conf"
	"Asura/src/logger"

	"github.com/Shopify/sarama"
	pkgerr "github.com/pkg/errors"
	cluster "github.com/bsm/sarama-cluster"
)

const (
	_family = "kafka"
)

var (
	// ErrProducer producer error.
	ErrProducer = errors.New("kafka producer nil")
	// ErrConsumer consumer error.
	ErrConsumer = errors.New("kafka consumer nil")
)

// Producer kafka.
type Producer struct {
	sarama.AsyncProducer
	sarama.SyncProducer
	c    *conf.KafkaProducer
	addr string
}

// NewProducer new kafka async producer and retry when has error.
func NewProducer(c *conf.KafkaProducer) (p *Producer) {
	var err error
	p = &Producer{
		c:    c,
		addr: strings.Join(c.Brokers, ","),
	}
	if !c.Sync {
		if err = p.asyncDial(); err != nil {
			go p.reAsyncDial()
		}
	} else {
		if err = p.syncDial(); err != nil {
			go p.reSyncDial()
		}
	}
	return
}

func (p *Producer) syncDial() (err error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_10_2_1
	p.SyncProducer, err = sarama.NewSyncProducer(p.c.Brokers, config)
	return
}

func (p *Producer) reSyncDial() {
	var err error
	for {
		if err = p.syncDial(); err == nil {
			logger.Info("kafka retry new sync producer ok")
			return
		}
		logger.Error("dial kafka producer error(%v)", err)
		time.Sleep(time.Second)
	}
}

func (p *Producer) asyncDial() (err error) {
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_2_1
	if p.AsyncProducer, err = sarama.NewAsyncProducer(p.c.Brokers, config); err == nil {
		go p.errproc()
		go p.successproc()
	}
	return
}

func (p *Producer) reAsyncDial() {
	var err error
	for {
		if err = p.asyncDial(); err == nil {
			logger.Info("kafka retry new async producer ok")
			return
		}
		logger.Error("dial kafka producer error(%v)", err)
		time.Sleep(time.Second)
	}
}

// errproc errors when aync producer publish messages.
// NOTE: Either Errors channel or Successes channel must be read. See the doc of AsyncProducer
func (p *Producer) errproc() {
	err := p.Errors()
	for {
		e, ok := <-err
		if !ok {
			return
		}
		logger.Error("kafka producer send message(%v) failed error(%v)", e.Msg, e.Err)
		if _, ok := e.Msg.Metadata.(context.Context); ok {
			// to do sth.
		}
	}
}

func (p *Producer) successproc() {
	suc := p.Successes()
	for {
		msg, ok := <-suc
		if !ok {
			return
		}
		if _, ok := msg.Metadata.(context.Context); ok {
			// to do sth.
		}
	}
}

// Input send msg to kafka
// NOTE: If producer has beed created failed, the message will lose.
func (p *Producer) Input(c context.Context, msg *sarama.ProducerMessage) (err error) {
	_, _ = msg.Key.Encode()
	if !p.c.Sync {
		if p.AsyncProducer == nil {
			err = ErrProducer
		} else {
			msg.Metadata = c
			p.AsyncProducer.Input() <- msg
		}
	} else {
		if p.SyncProducer == nil {
			err = ErrProducer
		} else {
			_, _, err = p.SyncProducer.SendMessage(msg)
		}
	}

	return pkgerr.WithStack(err)
}

// Close close producer.
func (p *Producer) Close() (err error) {
	if !p.c.Sync {
		if p.AsyncProducer != nil {
			return pkgerr.WithStack(p.AsyncProducer.Close())
		}
	}
	if p.SyncProducer != nil {
		return pkgerr.WithStack(p.SyncProducer.Close())
	}
	return
}


// Consumer kafka
type Consumer struct {
	ConsumerGroup *cluster.Consumer
	c             *conf.KafkaConsumer
}

// NewConsumer new a consumer.
func NewConsumer(c *conf.KafkaConsumer) (kc *Consumer) {
	var err error
	kc = &Consumer{
		c: c,
	}
	if c.Monitor != nil {
		go kc.monitor()
	}
	if err = kc.dial(); err != nil {
		go kc.redial()
	}
	return
}

func (c *Consumer) monitor() {
	mux := http.NewServeMux()
	mux.HandleFunc("/job/monitor/ping", ping)
	server := &http.Server{
		Addr:         c.c.Monitor.Addrs[0],
		Handler:      mux,
		ReadTimeout:  time.Duration(c.c.Monitor.ReadTimeout),
		WriteTimeout: time.Duration(c.c.Monitor.WriteTimeout),
	}
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server.ListenAndServe error(%v)", err)
		panic(err)
	}
}

func ping(wr http.ResponseWriter, r *http.Request) {
	return
}

func (c *Consumer) dial() (err error) {
	// cluster config
	cfg := cluster.NewConfig()
	// NOTE cluster auto commit offset interval
	cfg.Consumer.Offsets.CommitInterval = time.Second * 1
	// NOTE set fetch.wait.max.ms
	cfg.Consumer.MaxWaitTime = time.Millisecond * 100
	// NOTE errors that occur during offset management,if enabled, c.Errors channel must be read
	cfg.Consumer.Return.Errors = true
	// NOTE notifications that occur during consumer, if enabled, c.Notifications channel must be read
	cfg.Group.Return.Notifications = true
	cfg.Version = sarama.V0_10_2_1
	// The initial offset to use if no offset was previously committed.
	// default: OffsetOldest
	if strings.ToLower(c.c.Offset) != "new" {
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	// new cluster consumer
	c.ConsumerGroup, err = cluster.NewConsumer(c.c.Addrs, c.c.GroupID, c.c.Topic, cfg)
	return
}

func (c *Consumer) redial() {
	var err error
	for {
		if err = c.dial(); err == nil {
			logger.Info("kafka retry new consumer ok")
			return
		}
		logger.Error("dial kafka consumer error(%v)", err)
		time.Sleep(time.Second)
	}
}

// Close close consumer.
func (c *Consumer) Close() error {
	if c.ConsumerGroup != nil {
		return pkgerr.WithStack(c.ConsumerGroup.Close())
	}
	return nil
}
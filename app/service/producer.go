package service

import (
	"sync"
	"context"

	"Asura/conf"
	"Asura/app/queue/kafka"

	"github.com/Shopify/sarama"
)

// Pub define producer.
type Pub struct {
	group     string
	topic     string
	cluster   string
	appsecret string
	producer  *kafka.Producer
}

const (
	_group = "gukai-group"
	_topic = "gukai-test"
	_cluster = "test_kafka"
	_appsecret = "test"
)

var (
	// producer snapshot, key:group+topic
	producers = make(map[string]*kafka.Producer)
	pLock     sync.RWMutex
)

// NewPub new kafka producer.
func NewPub(c *conf.Config) (p *Pub) {
	p = &Pub{
		group: _group,
		topic: _topic,
		cluster: _cluster,
		appsecret: _appsecret,
		producer:  producer(_group, _topic, c.KafkaProducer.Test),
	}
	return
}

func producer(group, topic string, pCfg *conf.KafkaProducer) (p *kafka.Producer) {
	var (
		ok  bool
		key = genKey(group, topic)
	)
	pLock.RLock()
	if p, ok = producers[key]; ok {
		pLock.RUnlock()
		return
	}
	pLock.RUnlock()
	// new
	p = kafka.NewProducer(&conf.KafkaProducer{Brokers: pCfg.Brokers, Sync: pCfg.Sync, Cluster: _cluster})
	pLock.Lock()
	producers[key] = p
	pLock.Unlock()
	return
}

// Send publish kafka message.
func (p *Pub) Send(key, value []byte) (err error) {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	if err = p.producer.Input(context.TODO(), message); err != nil {
		return
	}
	return
}

func genKey(group, topic string) string {
	return group + ":" + topic
}

// Auth check user pub permission.
func (p *Pub) Auth(secret string) bool {
	return p.appsecret == secret
}
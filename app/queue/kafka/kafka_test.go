package kafka

import (
	"log"
	"sync"
	"testing"
	"time"

	"Asura/conf"
	"github.com/Shopify/sarama"
)

var testOne sync.Once

func TestKafka(t *testing.T) {
	testOne.Do(testConsumer)
	// producer
	cp := &conf.KafkaProducer{
		// Zookeeper: &conf.Zookeeper{Root: "/kafka", Addrs: []string{"172.16.33.54:2181"}, Timeout: time.Duration(time.Second)},
		Brokers:   []string{"172.16.33.54:9092"},
		Sync:      true,
	}
	p := NewProducer(cp)
	defer p.Close()
	msg := &sarama.ProducerMessage{
		Topic: "kafka_test",
		Key:   sarama.StringEncoder("test"),
		Value: sarama.StringEncoder("test"),
	}
	if partition, offset, err := p.SendMessage(msg); err != nil {
		t.Errorf("Kafka: send message error(%v)", err)
	} else {
		t.Logf("Kafka: partition: %d, offset: %d", partition, offset)
	}
}

func testConsumer() {
	// consumer
	cc := &conf.KafkaConsumer{
		Monitor:   &conf.HTTPServer{Addrs: []string{"0.0.0.0:8090"}},
		GroupID:     "test_group",
		Topic:    []string{"kafka_test"},
		Offset:    "old",
		Addrs:	  []string{"172.16.33.54:9092"},
		// Zookeeper: &conf.Zookeeper{Root: "/kafka", Addrs: []string{"172.16.33.54:2181"}, Timeout: xtime.Duration(time.Second)},
	}
	c := NewConsumer(cc)
	defer c.Close()
	time.Sleep(time.Second)
	go func() {
		for {
			if msg, ok := <-c.ConsumerGroup.Messages(); ok {
				log.Printf("Kafka: receive message: %v", msg)
			}
		}
	}()
	time.Sleep(time.Second * 3)
}

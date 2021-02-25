package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"Asura/app/dao"
	"Asura/conf"
	"Asura/app/queue/kafka"

	"github.com/Shopify/sarama"
)

const (
	_chanSize = 1024*1000
)

// Service abm service def.
type Service struct {
	c        *conf.Config
	dao	     *dao.Dao
	wg     	 *sync.WaitGroup
	quit  	 chan bool
	start    chan bool
	consumer *kafka.Consumer
	producer *Pub
	Commits  map[int32]*sarama.ConsumerMessage
	mapData  []Metadata
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		dao:  dao.New(c),
		wg:   new(sync.WaitGroup),
		quit: make(chan bool, 1),
		start: make(chan bool, 1),
		// consumer: NewSub(c),
		// producer: NewPub(c),
		Commits: make(map[int32]*sarama.ConsumerMessage),
	}

	// s.testProducer()
	// go s.testConsumer()
	return s
}

func (s *Service) Close() {
	s.quit <- true
	s.dao.Close()
}

func (s *Service) Wait() {
	time.Sleep(time.Millisecond)
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

func (s *Service) testProducer() {
	for i := 0; i <= 100; i++ {
		var key = []byte(UUid())
		var value = make(map[string]interface{})
		value["id"] = key
		value["name"] = "shies"
		str, _ := json.Marshal(value)
		fmt.Println(s.producer.Send(key, []byte(str)))
	}
}

func (s *Service) testConsumer() {
	for {
		var length int
		var err error
		length, err = s.IncrMessages(context.TODO())
		if err != nil {
			fmt.Println("The kafka transport redis fail(%v)", err)
			continue
		}
		fmt.Println(length, err)
		if length > 0 {
			s.ChunkTransport()
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *Service) ChunkTransport() {
	for {
		err := s.Transport(context.TODO())
		if err != nil {
			fmt.Println("The kafka transport redis fail(%v)", err)
			break
		}
		// 如果每1秒获取的增量数据过多，这里可以继续分批pipe打入redis
		// 最后标注kafka数据的偏移量，并且清空增量数据, 进行下次增量录入
		_ = s.Commit(context.TODO())
		break
	}
}
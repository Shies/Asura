package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"Asura/conf"
	"Asura/app/queue/kafka"
	log "Asura/src/logger"
)

const (
	_family     = "Ali-DTS"
	_actionSub  = "sub"
	_actionPub  = "pub"
	_actionAll  = "pubsub"
	_cmdPub     = "SET"
	_cmdSub     = "MGET"
	_authFormat = "%s:%s@%s/topic=%s&role=%s"
	_open       = int32(0)
	_closed     = int32(1)
)

var (
	// ErrAction action error.
	ErrAction = errors.New("action unknown")
	// ErrFull chan full
	ErrFull = errors.New("chan full")
	// ErrNoInstance no instances
	ErrNoInstance = errors.New("no Ali-DTS instances found")
)

// Message Data.
type Message struct {
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	Topic     string          `json:"topic"`
	Partition int32           `json:"partition"`
	Offset    int64           `json:"offset"`
	Timestamp int64           `json:"timestamp"`
}

func NewSub(c *conf.Config) (s *kafka.Consumer) {
	s = kafka.NewConsumer(c.KafkaConsumer.Test)
	return
}

func (s *Service) IncrMessages(c context.Context) (length int, err error) {
	ticker := time.NewTicker(time.Duration(time.Millisecond * time.Duration(600)))
Loop:
	for {
		select {
		case msg, ok := <-s.consumer.ConsumerGroup.Messages():
			if !ok {
				log.Error("Ali-DTS: %s binlog consumer exit!!!", s.c.KafkaConsumer.Test.GroupID)
				break Loop
			}
			s.Commits[msg.Partition] = msg
			rawValue := bytes.TrimFunc(msg.Value, func(r rune) bool {
				return r != 123 && r != 125 // 去除json垃圾byte在{}前后
			})
			for i, ch := range rawValue {
				switch {
				case ch > '~':   rawValue[i] = ' '
				case ch == '\r':
				case ch == '\n':
				case ch == '\t':
				case ch < ' ':   rawValue[i] = ' '
				}
			}
			parseMap, err := s.dao.JSON2map(rawValue)
			if err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(rawValue), err)
				continue
			}
			s.mapData = append(s.mapData, parseMap)
			/*
			// this is means when from dts read data then we need goto parse, e.g. parse fail
			var parseMap1, parseMap2 map[string]interface{}
			parseParts := strings.Split(string(rawValue), " string ")
			if len(parseParts) > 0 {
				parseMap1, err = s.dao.JSON2map([]byte(parseParts[0]))
				if err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", string(parseParts[0]), err)
					continue
				}
				parseMap2, err = s.dao.JSON2map([]byte(parseParts[1]))
				if err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", string(parseParts[1]), err)
					continue
				}
				s.mapData = append(s.mapData, parseMap2)
			} else {
				parseMap1, err = s.dao.JSON2map([]byte(parseParts[0]))
				if err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", string(parseParts[0]), err)
					continue
				}
			}
			s.mapData = append(s.mapData, parseMap1)
			*/
		case err := <-s.consumer.ConsumerGroup.Errors():
			log.Error("kafka:Error:%v", err.Error())
		case note := <-s.consumer.ConsumerGroup.Notifications():
			log.Info("kafka:Note: %v", note)
		case <-ticker.C:
			break Loop
		}
	}
	// todo: 额外的参数
	length = len(s.mapData)
	return
}

func (s *Service) Commit(c context.Context) (err error) {
	for k, msg := range s.Commits {
		s.consumer.ConsumerGroup.MarkOffset(msg, "")
		delete(s.Commits, k)
	}
	s.mapData = []Metadata{}
	return
}

func (s *Service) Transport(c context.Context) (err error) {
	var (
		values = make(map[string]string)
		oldValues = make(map[string]string)
		rawValue json.RawMessage
	)
	for _, msg := range s.mapData {
		typ := msg.Get("type")
		switch typ {
		case "UPDATE":
			if data, ok := msg.Get("data").([]interface{}); ok {
				for _, val := range data {
					meta := val.(map[string]interface{})
					if sn, ok := meta["sn"]; ok {
						if _, ok := values[sn.(string)]; !ok {
							rawValue, err = json.Marshal(meta)
							values[sn.(string)] = string(rawValue)
						}
					}
					break
				}
			}
			if data, ok := msg.Get("old").([]interface{}); ok {
				for _, val := range data {
					meta := val.(map[string]interface{})
					if sn, ok := meta["sn"]; ok {
						if _, ok := oldValues[sn.(string)]; !ok {
							rawValue, err = json.Marshal(meta)
							oldValues[sn.(string)] = string(rawValue)
						}
					}
					break
				}
			}
		case "INSERT":
		case "DELETE":
		default:
		}
	}
	if len(values) > 0 {
		fmt.Println(len(values)) // drop repeat after of result
		reply, err := s.dao.PipeLine(values)
		fmt.Println(reply, err) // the pipeline transport whether success
		if reply.(string) != "OK" {
			return err
		}
	}
	return
}

func (s *Service) AllMessages(c context.Context) (length int, err error) {
	return
}
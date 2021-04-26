package service

import (
	"context"
	"encoding/json"

	"Asura/app/model"
	log "Asura/src/logger"
)

func (s *Service) Test(ctx context.Context, key string) (test *model.Test, err error) {
	reply, err := s.dao.GetValues(key)
	if err != nil {
		log.Error("get values cache fail(%v)", err)
		return
	}

	test = &model.Test{}
	err = json.Unmarshal([]byte(reply), test)
	if err != nil {
		log.Error("convert json to struct fail(%v)", err)
		return
	}

	return
}

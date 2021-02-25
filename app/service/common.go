package service

import (
	"regexp"
	"math"
	"time"
	"crypto/md5"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

const (
	_regular = "^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\\\\d{8}$"
)

func UUid() string {
	v := uuid.Must(uuid.NewV4(), nil).String()
	has := md5.Sum([]byte(v))
	md5val := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5val
}

func ConvertTime(str string, layout string) time.Time {
	timer, _ := time.ParseInLocation(layout, str, time.Local)
	return timer
}

func IsMobile(mobileNum string) bool {
	reg := regexp.MustCompile(_regular)
	return reg.MatchString(mobileNum)
}

func InArray(source int64, target []int64) bool {
	var isExists = false
	for _, v := range target {
		if v == source {
			isExists = true
			break
		}
	}

	return isExists
}

func ValidateIsTimeout(timef int64) bool {
	timef = int64(math.Round(float64(timef) / 1000))
	if (time.Now().Unix() - timef) > 60 {
		return true
	}

	return false
}
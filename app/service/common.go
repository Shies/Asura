package services

import (
	"crypto/md5"
	"fmt"
	"math"
	"regexp"
	"time"

	xtime "Asura/app/time"

	uuid "github.com/satori/go.uuid"
)

const (
	_regular = "^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\\\\d{8}$"
)

func UUid() string {
	v := uuid.Must(uuid.NewV4(), nil).String()
	return Md5(v)
}

func Md5(value string) string {
	has := md5.Sum([]byte(value))
	md5val := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5val
}

func ConvertTime(str string, layout string) time.Time {
	dateTime, _ := time.ParseInLocation(layout, str, xtime.CSTZone)
	return dateTime
}

func IsMobile(mobile string) bool {
	reg := regexp.MustCompile(_regular)
	return reg.MatchString(mobile)
}

func InArray(source int64, target []int64) bool {
	var exists = false
	for _, v := range target {
		if v == source {
			exists = true
			break
		}
	}

	return exists
}

func ValidateIsTimeout(timef float64) bool {
	timeout := int64(math.Round(float64(timef)))
	// fmt.Println(time.Now().Unix())
	// fmt.Println(timeout)
	if (time.Now().Unix() - timeout) >= 0 {
		return true
	}

	return false
}

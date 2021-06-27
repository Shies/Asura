package service

import (
	"reflect"
	"strconv"
	"unsafe"

	"Asura/app/model"
	xtime "Asura/app/time"

	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

func (s *Service) map2struct(result interface{}, stype string) interface{} {
	var (
		refvalue reflect.Value
		reftype  reflect.Type
		source   interface{}
	)
	if realTypeResult, ok := result.(php_serialize.PhpArray); ok {
		if (stype == "welcome") {
			source = &model.Welcome{}
		} else {
			source = unsafe.Pointer(nil)
		}
		refvalue = reflect.ValueOf(source)
		reftype  = reflect.TypeOf(source)
		if reftype.Kind() != reflect.Ptr {
			panic("err: type invalid!")
		}
		for i := 0; i < reftype.Elem().NumField(); i++ {
			val := realTypeResult[reftype.Elem().Field(i).Tag.Get("json")]
			if val != nil {
				kind := reftype.Elem().Field(i).Type.Kind()
				switch kind {
				case reflect.Int64,reflect.Int,reflect.Int32,reflect.Int16,reflect.Int8:
					var origin int
					if str, ok := val.(string); ok {
						origin, _ = strconv.Atoi(str)
					} else {
						origin = int(val.(int64))
					}
					refvalue.Elem().Field(i).SetInt(int64(origin))
				case reflect.Float64,reflect.Float32:
					var origin int
					if str, ok := val.(string); ok {
						origin, _ = strconv.Atoi(str)
					} else {
						origin = int(val.(float64))
					}
					refvalue.Elem().Field(i).SetFloat(float64(origin))
				case reflect.String:
					refvalue.Elem().Field(i).SetString(val.(string))
				case reflect.Bool:
					var origin int
					if origin = int(val.(int64)); origin > 0 {
						refvalue.Elem().Field(i).SetBool(true)
					} else {
						refvalue.Elem().Field(i).SetBool(false)
					}
				case reflect.Struct:
					// 时间类型
					var t xtime.DateTime
					t = xtime.DateTime{Time: ConvertTime(val.(string), xtime.StandardLayout)}
					refvalue.Elem().Field(i).Set(reflect.ValueOf(t))
				}
			}
		}
	}

	return source
}
package errors

import (
	"encoding/json"
	"reflect"
	"runtime"
	"strconv"

	"github.com/khacman98/goutils/log"
)

// struct for detailed error construct contains code, message, type and additional data
type ErrorItem struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Caller  *CallerItem `json:"caller"`
	Func    string      `json:"func"`
}

type CallerItem struct {
	File string
	Line int
}

func (item *ErrorItem) Error() string {
	return "[" + strconv.Itoa(item.Code) + "]" + item.Message + " -> " + item.ToJSON()
}

func (item *ErrorItem) ToJSON() string {
	json, _ := json.Marshal(item)
	return string(json)
}

// constructor for new ErrorItem
func New(
	code int,
	message string,
	data interface{},
) error {

	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	err := ErrorItem{
		Code:    code,
		Message: message,
		Data:    data,
		Caller: &CallerItem{
			File: file,
			Line: line,
		},
		Func: f.Name(),
	}

	//エラー発生時にログも出す
	log.LogfWithDepth(log.LOG_LEVEL_ERROR, 2, "=====================")
	log.LogfWithDepth(log.LOG_LEVEL_ERROR, 2, "Error Occurs")
	log.LogfWithDepth(log.LOG_LEVEL_ERROR, 2, "File:%s", file)
	log.LogfWithDepth(log.LOG_LEVEL_ERROR, 2, "Line:%d", line)
	log.LogfWithDepth(log.LOG_LEVEL_ERROR, 2, "message:%s", message)
	log.LogfWithDepth(log.LOG_LEVEL_ERROR, 2, "message:%v", data)
	log.LogfWithDepth(log.LOG_LEVEL_ERROR, 2, "=====================")

	return &err
}

func IsErrorItem(err error) bool {
	if reflect.TypeOf(err).String() == "*errors.ErrorItem" {
		return true
	} else {
		return false
	}
}

func ToErrorItem(err error) *ErrorItem {
	if reflect.TypeOf(err).String() == "*errors.ErrorItem" {
		return err.(*ErrorItem)
	} else {
		return nil
	}
}

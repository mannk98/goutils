package response

import (
	"encoding/json"
	"net/http"
	"os"
)

var ERRORS = map[string]string{
	"E99999": "Something went wrong please try again.!",
	"E00001": "ClientId does not exists.",
}

type ApiResponseStatus struct {
	Code int    `json:"code"`
	Type string `json:"type"`
}

type ApiResponse struct {
	ProcessId string            `json:"process_id"`
	Path      string            `json:"path"`
	Status    ApiResponseStatus `json:"status"`
	Request   interface{}       `json:"request"`
	Errors    []error           `json:"errors"`
	Data      interface{}       `json:"data"`
}

// convert to json string
func (item *ApiResponse) ToJSON() string {
	jsonout, err := json.Marshal(item)
	if err != nil {
		return ""
	}
	return string(jsonout)
}

func NewApiResponse(path string) *ApiResponse {
	resp_status := ApiResponseStatus{
		Code: http.StatusOK,
		Type: http.StatusText(http.StatusOK),
	}
	resp_error := []error{}
	resp_data := make([]interface{}, 0)

	return &ApiResponse{
		ProcessId: os.Getenv("PROCESS_ID"),
		Status:    resp_status,
		Errors:    resp_error,
		Data:      resp_data,
		Path:      path,
	}
}

// Response for [Get] API
type ApiGetResponse struct {
	ProcessId    string            `json:"process_id"`
	Path         string            `json:"path"`
	Status       ApiResponseStatus `json:"status"`
	Request      interface{}       `json:"request"`
	Errors       []error           `json:"errors"`
	ContentRange string            `json:"content_range"`
	Data         interface{}       `json:"data"`
}

// convert to json string
func (item *ApiGetResponse) ToJSON() string {
	jsOut, err := json.Marshal(item)
	if err != nil {
		return ""
	}
	return string(jsOut)
}

func NewApiGetResponse(path string) *ApiGetResponse {
	respStatus := ApiResponseStatus{
		Code: http.StatusOK,
		Type: http.StatusText(http.StatusOK),
	}
	respError := []error{}
	respData := make([]interface{}, 0)

	return &ApiGetResponse{
		ProcessId: os.Getenv("PROCESS_ID"),
		Status:    respStatus,
		Errors:    respError,
		Data:      respData,
		Path:      path,
	}
}

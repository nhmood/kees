package helpers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"kees/server/messages"
)

func Format(data interface{}) ([]byte, error) {
	output, err := json.Marshal(data)
	return output, err
}

func Parse(r *http.Request, data interface{}) error {
	a := json.NewDecoder(r.Body)
	a.Decode(data)
	return nil
}

func QueryInt(r *http.Request, p string) int64 {
	q := QueryString(r, p)
	param, _ := strconv.ParseInt(q, 10, 64)
	return param
}

func QueryString(r *http.Request, p string) string {
	params := r.URL.Query()
	param := params.Get(p)
	return param
}

func GetPage(r *http.Request) int64 {
	queryParams := r.URL.Query()
	page, _ := strconv.ParseInt(queryParams.Get("page"), 10, 64)
	Debug(page)

	return page
}

type Location int8

const (
	URLParam Location = iota
	QueryParam
)

func GetIntParam(r *http.Request, path string, location Location) int64 {
	var param string
	switch location {
	case URLParam:
		vars := mux.Vars(r)
		param = vars[path]
	case QueryParam:
		queryParams := r.URL.Query()
		param = queryParams.Get(path)
	}

	paramID, _ := strconv.ParseInt(param, 10, 64)
	return paramID
}

func GetStringParam(r *http.Request, path string, location Location) string {
	var param string
	switch location {
	case URLParam:
		vars := mux.Vars(r)
		param = vars[path]
	case QueryParam:
		queryParams := r.URL.Query()
		param = queryParams.Get(path)
	}

	return param
}

func ToInterface(data interface{}) map[string]interface{} {
	str, _ := Format(data)
	var conv map[string]interface{}
	json.Unmarshal([]byte(str), &conv)
	return conv
}

func ToStruct(data interface{}, target interface{}) {
	str, _ := Format(data)
	json.Unmarshal([]byte(str), target)
}

func Halt(w http.ResponseWriter, code int, message string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}

	response := messages.Generic{
		Message: message,
		Data:    data,
	}

	payload, err := Format(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(payload)
	return
}

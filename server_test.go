package main

import (
	"avitoChatAPI/requestTypes"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TestUserReq struct {
	name       string
	method     string
	input      *requestTypes.UserRequest
	want       string
	statusCode int
}

type TestUserChatsReq struct {
	name       string
	method     string
	input      *requestTypes.UserChatsRequest
	want       string
	statusCode int
}

type TestMsgReq struct {
	name       string
	method     string
	input      *requestTypes.MessageRequest
	want       string
	statusCode int
}

type TestChatReq struct {
	name       string
	method     string
	input      *requestTypes.ChatRequest
	want       string
	statusCode int
}

func handleTest(serverHandler func(w http.ResponseWriter, r *http.Request), tests []TestUserReq, t *testing.T, url string) {
	handler := http.HandlerFunc(serverHandler)
	for _, tr := range tests {
		t.Run(tr.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			bodyBytes, _ := json.Marshal(tr.input)
			body := bytes.NewReader(bodyBytes)
			req, _ := http.NewRequest(tr.method, url, body)
			handler(rec, req)
			data := rec.Body.String()
			data = data[0 : len(data)-1]
			if data != tr.want {
				t.Errorf("Got %s Expected %s", data, tr.want)
			}
		})
		fmt.Println(tr.name, "OK!")
	}
}

func TestAddUser(t *testing.T) {
	testReq := []TestUserReq{
		{
			name:       "normal user",
			method:     http.MethodPost,
			input:      &requestTypes.UserRequest{Name: "qwerty"},
			want:       `{"id":1}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "error user",
			method:     http.MethodPost,
			input:      &requestTypes.UserRequest{},
			want:       `400`,
			statusCode: http.StatusOK,
		},
	}
	dbPointer, _ = initDataBase("./tmpAvito.db")
	handleTest(addUser, testReq, t, "/users/add")
	os.Remove("./tmpAvito.db")
}

package main

import (
	"avitoChatAPI/requestTypes"
	"avitoChatAPI/responseTypes"
	"avitoChatAPI/sqlScheme"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var dbPointer *sql.DB

//Сообщение об ошибке
func respondError(w http.ResponseWriter, err error, status int) {
	log.Println(err.Error(), status)
	http.Error(w, strconv.Itoa(status), status)
}

//Отправка ответа
func sendRespond(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//Парсинг JSON-а
func decode(decodeStruct interface{}, requestBody io.ReadCloser) error {
	decoder := json.NewDecoder(requestBody)
	decoder.DisallowUnknownFields()
	return decoder.Decode(decodeStruct)
}

//Обработка /users/add
func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser requestTypes.UserRequest
	createdAt := time.Now().Format(time.RFC3339)

	if err := decode(&newUser, r.Body); err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}

	if len(newUser.Name) == 0 {
		respondError(w, errors.New("Empty Name!"), http.StatusBadRequest)
		return
	}

	ret, err := dbPointer.Exec("INSERT INTO users (name, createdAt) VALUES (?, ?);", newUser.Name, createdAt)

	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	lastId, _ := ret.LastInsertId()

	response := responseTypes.IdResponse{lastId}

	sendRespond(w, response)
}

//Обработка /chats/add
func addChat(w http.ResponseWriter, r *http.Request) {

	var newChat requestTypes.ChatRequest
	createdAt := time.Now().Format(time.RFC3339)

	if err := decode(&newChat, r.Body); err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}

	if len(newChat.Name) == 0 {
		respondError(w, errors.New("Empty Name!"), http.StatusBadRequest)
		return
	}

	if len(newChat.Users) == 0 {
		respondError(w, errors.New("Empty Users!"), http.StatusBadRequest)
		return
	}

	tx, err := dbPointer.Begin()

	ret, err := tx.Exec("INSERT INTO chats (name, createdAt) VALUES (?, ?);", newChat.Name, createdAt)

	if err != nil {
		tx.Rollback()
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	lastId, _ := ret.LastInsertId()

	for _, v := range newChat.Users {
		userId := v
		_, err := tx.Exec("INSERT INTO chatsJoinUsers (chat_id, user_id) VALUES (?, ?);", lastId, userId)
		if err != nil {
			tx.Rollback()
			respondError(w, err, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	response := responseTypes.IdResponse{lastId}

	sendRespond(w, response)
}

//Обработка /messages/add
func addMessage(w http.ResponseWriter, r *http.Request) {
	var newMessage requestTypes.MessageRequest
	createdAt := time.Now().Format(time.RFC3339)

	if err := decode(&newMessage, r.Body); err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}

	ret, err := dbPointer.Exec("INSERT INTO messages (text, chat_id, user_id, createdAt) VALUES (?, ?, ?, ?);", newMessage.Text, newMessage.ChatId, newMessage.UserId, createdAt)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}
	lastId, _ := ret.LastInsertId()

	response := responseTypes.IdResponse{lastId}

	sendRespond(w, response)
}

//Обработка /chats/get
func getUserChats(w http.ResponseWriter, r *http.Request) {
	var userChats requestTypes.UserChatsRequest

	if err := decode(&userChats, r.Body); err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}

	userId := userChats.UserId

	var userCheck int

	selectErr := dbPointer.QueryRow("SELECT id FROM users where id=?", userId).Scan(&userCheck)
	if selectErr != nil {
		if selectErr != sql.ErrNoRows {
			respondError(w, selectErr, http.StatusBadRequest)
			return
		}
		respondError(w, errors.New("No such user!"), http.StatusBadRequest)
		return
	}

	rows, err := dbPointer.Query("SELECT * FROM chats where id IN (SELECT chat_id FROM chatsJoinUsers WHERE user_id = ?) ORDER BY createdAt DESC", userId)

	defer rows.Close()

	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	response := make([]responseTypes.ChatResponse, 0, 32)

	for rows.Next() {
		var chatResponse responseTypes.ChatResponse
		err := rows.Scan(&chatResponse.Id, &chatResponse.Name, &chatResponse.CreatedAt)
		if err != nil {
			respondError(w, err, http.StatusInternalServerError)
			return
		}
		response = append(response, chatResponse)
	}

	sendRespond(w, response)
}

//Обработка /messages/get
func getChatMessages(w http.ResponseWriter, r *http.Request) {
	var chat requestTypes.ChatMessagesRequest
	if err := decode(&chat, r.Body); err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}
	chatId := chat.ChatId

	var chatCheck int

	selectErr := dbPointer.QueryRow("SELECT id FROM chats where id=?", chatId).Scan(&chatCheck)

	if selectErr != nil {
		if selectErr != sql.ErrNoRows {
			respondError(w, selectErr, http.StatusBadRequest)
			return
		}
		respondError(w, errors.New("No such chat!"), http.StatusBadRequest)
		return
	}

	rows, err := dbPointer.Query("SELECT * FROM messages where chat_id=? ORDER BY createdAt DESC", chatId)

	defer rows.Close()

	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	response := make([]responseTypes.MessageResponse, 0, 32)

	for rows.Next() {
		var messageResp responseTypes.MessageResponse
		err := rows.Scan(&messageResp.Id, &messageResp.Text, &messageResp.ChatId, &messageResp.UserId, &messageResp.CreatedAt)
		if err != nil {
			respondError(w, err, http.StatusInternalServerError)
			return
		}
		response = append(response, messageResp)
	}

	sendRespond(w, response)
}

//Инициализация базы данных
func initDataBase(name string) (*sql.DB, error) {

	database, connectionError := sql.Open("sqlite3", name)

	if connectionError != nil {
		return nil, connectionError
	}

	database.Exec("PRAGMA foreign_keys = ON")

	_, err := database.Exec(sqlScheme.Scheme)
	if err != nil {
		return nil, err
	}

	fmt.Println("Init DB success")
	return database, nil
}

const (
	port = ":9000"
)

func main() {
	var errDb error
	dbPointer, errDb = initDataBase("./avitoChat.db")
	if errDb != nil {
		log.Fatal("Can't init database!\n", errDb)
		return
	}

	defer dbPointer.Close()

	fmt.Println("Started server..")
	http.HandleFunc("/users/add", addUser)
	http.HandleFunc("/chats/add", addChat)
	http.HandleFunc("/messages/add", addMessage)
	http.HandleFunc("/chats/get", getUserChats)
	http.HandleFunc("/messages/get", getChatMessages)

	log.Fatal(http.ListenAndServe(port, nil))
}

package requestTypes

//Типы для http-запросов

type ChatRequest struct {
	Name  string  `json:"name"`
	Users []int64 `json:"users"`
}

type MessageRequest struct {
	ChatId int64  `json:"chat"`
	UserId int64  `json:"author"`
	Text   string `json:"text"`
}

type UserChatsRequest struct {
	UserId int64 `json:"user"`
}

type UserRequest struct {
	Name string `json:"username"`
}

type ChatMessagesRequest struct {
	ChatId int64 `json:"chat"`
}

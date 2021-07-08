package responseTypes

//Типы для http-ответов

type MessageResponse struct {
	Id        int64  `json:"id"`
	Text      string `json:"text"`
	ChatId    int64  `json:"chat_id"`
	UserId    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type ChatResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type IdResponse struct {
	Id int64 `json:"id"`
}

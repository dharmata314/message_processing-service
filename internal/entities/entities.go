package entities

type Message struct {
	ID      int    `json:"message_id"`
	Content string `json:"message_content"`
	Status  string `json:"message_status"`
}

type User struct {
	ID       int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

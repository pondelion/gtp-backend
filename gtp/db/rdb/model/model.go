package model

type NewTodo struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type Todo struct {
	ID     int    `json:"id" gorm:"AUTO_INCREMENT"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	UserID int
	User   *User `json:"user" gorm:"ForeignKey:UserID"`
}

type User struct {
	ID   int    `json:"id" gorm:"AUTO_INCREMENT"`
	Name string `json:"name"`
}

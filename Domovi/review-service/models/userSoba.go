package models

type User struct {
	ID      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
}

type Soba struct {
	ID         string `json:"id" bson:"_id"`
	RoomNumber string `json:"roomNumber" bson:"roomNumber"`
}

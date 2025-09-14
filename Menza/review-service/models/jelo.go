package models

type User struct {
	ID string `json:"id" bson:"_id"`
}

type Jelo struct {
	ID string `json:"jeloId" bson:"jeloId"`
}

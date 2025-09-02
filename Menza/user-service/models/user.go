package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Password       string             `bson:"password" json:"password"`
	Role           string             `bson:"role" json:"role"`
	Name           string             `bson:"name" json:"name"`
	Surname        string             `bson:"surname" json:"surname"`
	Email          string             `bson:"email" json:"email"`
	IsActive       bool               `bson:"isActive" json:"isActive"`
	Alergije       []string           `bson:"alergije,omitempty" json:"alergije,omitempty"`
	OmiljenaJela   []string           `bson:"omiljenaJela,omitempty" json:"omiljenaJela,omitempty"`
	IndeksStudenta string             `bson:"indeksStudenta,omitempty" json:"indeksStudenta,omitempty"`
}

func NewUser(username, password, role, name, surname, email string) User {
	return User{
		Password:       password,
		Role:           role,
		Name:           name,
		Surname:        surname,
		Email:          email,
		IsActive:       false,
		Alergije:       []string{}, // prazna lista po defaultu
		OmiljenaJela:   []string{}, // prazna lista po defaultu
		IndeksStudenta: "",         // prazno po defaultu
	}
}

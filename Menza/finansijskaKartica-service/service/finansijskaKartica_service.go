package service

import (
	"context"
	"finansijskaKartica-service/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FinansijskaKarticaService struct {
	Collection *mongo.Collection
}

// Konstruktor: prima bazu koja je već inicijalizovana
func NewFinansijskaKarticaService(db *mongo.Database) *FinansijskaKarticaService {
	return &FinansijskaKarticaService{
		Collection: db.Collection("finansijske_kartice"),
	}
}

// Kreira novu karticu
func (s *FinansijskaKarticaService) CreateKartica(kartica models.FinansijskaKartica) (models.FinansijskaKartica, error) {
	_, err := s.Collection.InsertOne(context.TODO(), kartica)
	if err != nil {
		return models.FinansijskaKartica{}, err
	}
	return kartica, nil
}

// Vraća sve kartice
func (s *FinansijskaKarticaService) GetKartice() ([]models.FinansijskaKartica, error) {
	var kartice []models.FinansijskaKartica
	cursor, err := s.Collection.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &kartice); err != nil {
		return nil, err
	}
	return kartice, nil
}

// NOVO: Vraća karticu po userId
func (s *FinansijskaKarticaService) GetKarticaByUserID(userID primitive.ObjectID) (models.FinansijskaKartica, error) {
	var kartica models.FinansijskaKartica
	err := s.Collection.FindOne(context.TODO(), bson.M{"userId": userID}).Decode(&kartica)
	if err == mongo.ErrNoDocuments {
		return models.FinansijskaKartica{}, fmt.Errorf("kartica for user %s not found", userID.Hex())
	} else if err != nil {
		return models.FinansijskaKartica{}, err
	}
	return kartica, nil
}

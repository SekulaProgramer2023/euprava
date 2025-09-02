package bootstrap

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"jelovnik-service/db"
	"jelovnik-service/models"
	"os"
)

func InsertInitialJela() {
	if os.Getenv("ENABLE_BOOTSTRAP") != "true" {
		return
	}

	collection := db.Client.Database("eupravaM").Collection("jela")

	// Proveri da li kolekcija već ima jela
	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error counting jela:", err)
		return
	}

	if count > 0 {
		return // Preskoči ako već postoje
	}

	var jela []interface{}

	jela = append(jela, models.Jelo{
		JeloID:     primitive.NewObjectID(),
		Naziv:      "Ćevapi sa kajmakom",
		Kategorija: models.Meso,
		TipObroka:  models.Rucak,
		Kalorije:   550,
		Nutritijenti: map[string]float64{
			"proteini":       35,
			"masti":          30,
			"ugljeniHidrati": 40,
		},
	})

	jela = append(jela, models.Jelo{
		JeloID:     primitive.NewObjectID(),
		Naziv:      "Palačinke sa džemom",
		Kategorija: models.Desert,
		TipObroka:  models.Dorucak,
		Kalorije:   350,
		Nutritijenti: map[string]float64{
			"proteini":       8,
			"masti":          10,
			"ugljeniHidrati": 60,
		},
	})

	jela = append(jela, models.Jelo{
		JeloID:     primitive.NewObjectID(),
		Naziv:      "Salata od povrća",
		Kategorija: models.Salata,
		TipObroka:  models.Rucak,
		Kalorije:   150,
		Nutritijenti: map[string]float64{
			"proteini":       5,
			"masti":          2,
			"ugljeniHidrati": 20,
		},
	})

	jela = append(jela, models.Jelo{
		JeloID:     primitive.NewObjectID(),
		Naziv:      "Čorba od povrća",
		Kategorija: models.Kuvano,
		TipObroka:  models.Rucak,
		Kalorije:   200,
		Nutritijenti: map[string]float64{
			"proteini":       6,
			"masti":          3,
			"ugljeniHidrati": 25,
		},
	})

	jela = append(jela, models.Jelo{
		JeloID:     primitive.NewObjectID(),
		Naziv:      "Omlet sa povrćem",
		Kategorija: models.Kuvano,
		TipObroka:  models.Dorucak,
		Kalorije:   300,
		Nutritijenti: map[string]float64{
			"proteini":       20,
			"masti":          15,
			"ugljeniHidrati": 10,
		},
	})

	jela = append(jela, models.Jelo{
		JeloID:     primitive.NewObjectID(),
		Naziv:      "Pasta sa povrćem",
		Kategorija: models.Vegetarijansko,
		TipObroka:  models.Rucak,
		Kalorije:   450,
		Nutritijenti: map[string]float64{
			"proteini":       12,
			"masti":          8,
			"ugljeniHidrati": 80,
		},
	})

	jela = append(jela, models.Jelo{
		JeloID:     primitive.NewObjectID(),
		Naziv:      "Voćna salata",
		Kategorija: models.Desert,
		TipObroka:  models.Dorucak,
		Kalorije:   200,
		Nutritijenti: map[string]float64{
			"proteini":       2,
			"masti":          1,
			"ugljeniHidrati": 50,
		},
	})

	jela = append(jela, models.Jelo{
		JeloID:     primitive.NewObjectID(),
		Naziv:      "Pečena piletina sa povrćem",
		Kategorija: models.Meso,
		TipObroka:  models.Vecera,
		Kalorije:   600,
		Nutritijenti: map[string]float64{
			"proteini":       40,
			"masti":          25,
			"ugljeniHidrati": 30,
		},
	})

	// Ubaci u bazu
	_, err = collection.InsertMany(context.TODO(), jela)
	if err != nil {
		fmt.Println("Error inserting initial jela:", err)
	} else {
		fmt.Println("Inserted initial jela into database")
	}
}

func ClearJela() {
	collection := db.Client.Database("eupravaM").Collection("jela")
	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error clearing jela:", err)
	} else {
		fmt.Println("Cleared jela from database")
	}
}

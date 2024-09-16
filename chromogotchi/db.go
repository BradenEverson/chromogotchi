package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE   = "chromogotchi"
	COLLECTION = "pets"
)

var client *mongo.Client

func connectClient() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("MONGODB_URI Not Set >:((((")
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}

	client = mongoClient

	return nil
}

func loadPetsFromMongo() error {
	if client == nil {
		return errors.New("Client has not yet connected")
	}

    collection := client.Database(DATABASE).Collection(COLLECTION)
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{}, options.Find())
    if err != nil {
        return errors.New("Error retrieving pets from MongoDB")
    }
    defer cursor.Close(ctx)

    allPets = make(map[string]Pet)

    for cursor.Next(ctx) {
        var pet Pet
        if err := cursor.Decode(&pet); err != nil {
            log.Println("Error decoding pet document:", err)
            continue
        }
        
        allPets[pet.Id] = pet
    }

    if err := cursor.Err(); err != nil {
        return errors.New("Error iterating")
    }

    log.Println("Loaded pets from MongoDB:", allPets)
    return nil
}

func savePetsToMongo() {
	if client == nil {
		log.Fatal("Client has not yet connected")
	}

	for {
		time.Sleep(30 * time.Minute)
		fmt.Println("Saving all pets to MongoDB...")

		collection := client.Database(DATABASE).Collection(COLLECTION)
		for _, pet := range allPets {
			filter := bson.M{"id": pet.Id}
            update := bson.M{
                "$set": bson.M{
                    "name":        pet.Name,
                    "hunger":      pet.Hunger,
                    "happiness":   pet.Happiness,
                    "sleepiness":  pet.Wakefullness,
                    "hungerRate":  pet.HungerRate,
                    "depression":  pet.Depression,
                    "sleepyRate":  pet.SleepyRate,
                },
            }
            _, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
			if err != nil {
				log.Printf("Error saving pet %s: %v\n", pet.Id, err)
			}
		}
	}
}

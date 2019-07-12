package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://172.28.18.24:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	collection := client.Database("learn").Collection("unicorns")
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(result)
		}
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
}

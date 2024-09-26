package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoggerRepository struct {
	client *mongo.Client
}

func NewLoggerRepository(mongoClient *mongo.Client) *LoggerRepository {
	return &LoggerRepository{
		client: mongoClient,
	}
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (repo *LoggerRepository) Insert(logEntry LogEntry) error {
	collection := repo.client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      logEntry.Name,
		Data:      logEntry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs", err)
		return err
	}
	return nil
}

func (repo *LoggerRepository) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := repo.client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("error decoding log into slice:", err)
		} else {
			logs = append(logs, &item)
		}
	}
	return logs, nil
}

func (repo *LoggerRepository) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := repo.client.Database("logs").Collection("logs")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var logEntry LogEntry
	singleResult := collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&logEntry)
	if singleResult.Error() != "" {
		return nil, err
	}

	return &logEntry, nil
}

func (repo *LoggerRepository) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := repo.client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (repo *LoggerRepository) Update(log LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := repo.client.Database("logs").Collection("logs")

	objectId, err := primitive.ObjectIDFromHex(log.ID)
	if err != nil {
		return nil, err
	}

	updateResult, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: log.Name},
				{Key: "data", Value: log.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

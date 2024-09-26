package data

import "go.mongodb.org/mongo-driver/mongo"

type Repository interface {
	Insert(logEntry LogEntry) error
	All() ([]*LogEntry, error)
	GetOne(id string) (*LogEntry, error)
	DropCollection() error
	Update(log LogEntry) (*mongo.UpdateResult, error)
}

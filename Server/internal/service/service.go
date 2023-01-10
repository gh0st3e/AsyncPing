package service

import (
	"context"
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db    *sql.DB
	mongo *mongo.Client
}

func NewService(db *sql.DB, mongo *mongo.Client) *Service {
	return &Service{
		db:    db,
		mongo: mongo,
	}
}

func (s *Service) PingPSQL() error {
	err := s.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) PingMongo() error {
	err := s.mongo.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	return nil
}

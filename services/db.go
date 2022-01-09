package services

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBService interface {
	GetDB() *sql.DB
	CloseDB()
	InitSQL(sql string)
}

type Service struct {
	db *sql.DB
}

func NewDBService() DBService {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=postgres sslmode=disable")

	if err != nil {
		panic(err)
	}

	return Service{db}
}

func (s Service) GetDB() *sql.DB {
	return s.db
}

func (s Service) CloseDB() {
	s.db.Close()
}

func (s Service) InitSQL(sql string) {
	_, err := s.db.Exec(sql)

	if err != nil {
		panic(err)
	}
}

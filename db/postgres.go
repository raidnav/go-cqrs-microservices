package db

import (
	"context"
	"database/sql"
	"github.com/raidnav/go-cqrs-microservices/schema"
	"log"
)

//Repository interface is to handle interaction between db. It supports querying and inserting data to postgres.
type PostgresRepository struct {
	db *sql.DB
}

// Initiates connection to database.
func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

// Close connection to database
func (r *PostgresRepository) Close() {
	err := r.db.Close()
	if err != nil {
		log.Printf("Unable to close database connection")
	}
}

// InsertMeows is to append meow object to database
func (r *PostgresRepository) InsertMeows(ctx context.Context, meow schema.Meow) error {
	_, err := r.db.Exec("INSERT INTO meows(id, body, created_at) VALUES ($1, $2, $3)", meow.ID, meow.Body, meow.CreatedAt)
	return err
}

// ListMeows is to fetch all meows from database
func (r *PostgresRepository) ListMeows(ctx context.Context, skip uint64, take uint64) ([]schema.Meow, error) {
	rows, err := r.db.Query("SELECT * FROM meows ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse all rows into an array of Meows
	var meows []schema.Meow
	for rows.Next() {
		meow := schema.Meow{}
		if err = rows.Scan(&meow.ID, &meow.Body, &meow.CreatedAt); err == nil {
			meows = append(meows, meow)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return meows, nil
}

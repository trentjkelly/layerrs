package databaseRepository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthDatabaseRepository struct {
	db	*pgxpool.Pool
}

func NewAuthDatabaseRepository(db *pgxpool.Pool) *AuthDatabaseRepository {
	return &AuthDatabaseRepository{
		db: db,
	}
}

func (r *AuthDatabaseRepository) CloseDB() {
	r.db.Close()
}

// Creates a new magic link token for an artist
func (r *AuthDatabaseRepository) CreateMagicLinkToken(ctx context.Context, token string, artistId int) error {
	query := `INSERT INTO magic_link_tokens (hashed_token, artist_id) VALUES ($1, $2) RETURNING id;`
	row := r.db.QueryRow(ctx, query, token, artistId)
	
	var magicLinkTokenId int
	err := row.Scan(&magicLinkTokenId)
	if err != nil {
		return fmt.Errorf("failed to create magic link token in database: %w", err)
	}
	return nil
}
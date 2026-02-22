package databaseRepository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
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

func (r *AuthDatabaseRepository) GetMagicLinkToken(ctx context.Context, token string) (entities.MagicLinkToken, error) {
	query := `SELECT id, artist_id, created_at FROM magic_link_tokens WHERE hashed_token = $1;`
	row := r.db.QueryRow(ctx, query, token)

	var magicLinkToken entities.MagicLinkToken
	err := row.Scan(&magicLinkToken.Id, &magicLinkToken.ArtistId, &magicLinkToken.CreatedAt)
	if err != nil {
		return magicLinkToken, fmt.Errorf("failed to get magic link token from database: %w", err)
	}
	return magicLinkToken, nil
}
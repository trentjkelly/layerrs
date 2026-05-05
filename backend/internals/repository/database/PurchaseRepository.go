package databaseRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type PurchaseRepository struct {
	db *pgxpool.Pool
}

func NewPurchaseRepository(db *pgxpool.Pool) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) CloseDB() {
	r.db.Close()
}

func (r *PurchaseRepository) CreatePurchase(ctx context.Context, purchase *entities.Purchase) error {
	query := `INSERT INTO purchase (buyer_id, track_id, stripe_payment_intent_id, amount_cents) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	row := r.db.QueryRow(ctx, query, purchase.BuyerId, purchase.TrackId, purchase.StripePaymentIntentId, purchase.AmountCents)

	err := row.Scan(&purchase.Id, &purchase.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create purchase: %w", err)
	}

	return nil
}

func (r *PurchaseRepository) ReadPurchaseByPaymentIntentId(ctx context.Context, paymentIntentId string) (*entities.Purchase, error) {
	query := `SELECT id, buyer_id, track_id, stripe_payment_intent_id, amount_cents, created_at, updated_at FROM purchase WHERE stripe_payment_intent_id=$1`
	row := r.db.QueryRow(ctx, query, paymentIntentId)

	var purchase entities.Purchase
	err := row.Scan(&purchase.Id, &purchase.BuyerId, &purchase.TrackId, &purchase.StripePaymentIntentId, &purchase.AmountCents, &purchase.CreatedAt, &purchase.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read purchase by payment intent: %w", err)
	}

	return &purchase, nil
}

func (r *PurchaseRepository) ReadPurchasesByBuyerId(ctx context.Context, buyerId int) ([]entities.Purchase, error) {
	query := `SELECT id, buyer_id, track_id, stripe_payment_intent_id, amount_cents, created_at, updated_at FROM purchase WHERE buyer_id=$1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, buyerId)
	if err != nil {
		return nil, fmt.Errorf("failed to query purchases: %w", err)
	}
	defer rows.Close()

	var purchases []entities.Purchase
	for rows.Next() {
		var purchase entities.Purchase
		err = rows.Scan(&purchase.Id, &purchase.BuyerId, &purchase.TrackId, &purchase.StripePaymentIntentId, &purchase.AmountCents, &purchase.CreatedAt, &purchase.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase: %w", err)
		}
		purchases = append(purchases, purchase)
	}

	return purchases, nil
}

func (r *PurchaseRepository) HasPurchased(ctx context.Context, buyerId int, trackId int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM purchase WHERE buyer_id=$1 AND track_id=$2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, buyerId, trackId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check purchase: %w", err)
	}

	return exists, nil
}

func (r *PurchaseRepository) ReadTrackWithSeller(ctx context.Context, trackId int) (*entities.TrackWithSeller, error) {
	query := `
		SELECT t.id, t.artist_id, t.wav_r2_track_key, t.price_in_cents, a.stripe_account_id, a.stripe_account_status
		FROM track t
		JOIN artist a ON t.artist_id = a.id
		WHERE t.id = $1 AND t.is_valid = true
	`

	var trackWithSeller entities.TrackWithSeller
	var wavR2TrackKey sql.NullString
	var stripeAccountId sql.NullString

	err := r.db.QueryRow(ctx, query, trackId).Scan(
		&trackWithSeller.TrackId,
		&trackWithSeller.ArtistId,
		&wavR2TrackKey,
		&trackWithSeller.PriceInCents,
		&stripeAccountId,
		&trackWithSeller.StripeAccountStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read track with seller: %w", err)
	}

	if wavR2TrackKey.Valid {
		trackWithSeller.WavR2TrackKey = wavR2TrackKey.String
	}
	if stripeAccountId.Valid {
		trackWithSeller.StripeAccountId = stripeAccountId.String
	}

	return &trackWithSeller, nil
}

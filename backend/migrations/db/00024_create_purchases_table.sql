-- +goose Up
-- +goose StatementBegin
CREATE TABLE purchase (
    id SERIAL PRIMARY KEY,
    buyer_id INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE,
    track_id INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    stripe_payment_intent_id VARCHAR(255) NOT NULL,
    amount_cents INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(stripe_payment_intent_id)
);

CREATE INDEX idx_purchase_buyer_id ON purchase(buyer_id);
CREATE INDEX idx_purchase_track_id ON purchase(track_id);
CREATE INDEX idx_purchase_stripe_payment_intent_id ON purchase(stripe_payment_intent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS purchase;
-- +goose StatementEnd

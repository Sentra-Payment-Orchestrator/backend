CREATE TABLE IF NOT EXISTS payment_notifications (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id) ON DELETE CASCADE,
    raw_notification JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_notifications_transaction ON payment_notifications(transaction_id);
CREATE INDEX IF NOT EXISTS idx_payment_notifications_created ON payment_notifications(created_at);
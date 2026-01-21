CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    outlet_id INTEGER NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    transaction_number VARCHAR(100) UNIQUE NOT NULL,
    
    customer_name VARCHAR(255),
    customer_phone VARCHAR(50),
    customer_email VARCHAR(255),
    
    subtotal DECIMAL(15, 2) NOT NULL,
    tax_amount DECIMAL(15, 2) DEFAULT 0,
    discount_amount DECIMAL(15, 2) DEFAULT 0,
    total_amount DECIMAL(15, 2) NOT NULL,
    
    payment_method VARCHAR(50),
    payment_status VARCHAR(50) DEFAULT 'pending',
    payment_gateway_id INTEGER REFERENCES payment_gateways(id),
    
    midtrans_order_id VARCHAR(255),
    midtrans_transaction_id VARCHAR(255),
    midtrans_transaction_status VARCHAR(50),
    midtrans_fraud_status VARCHAR(50),
    midtrans_payment_type VARCHAR(50),
    midtrans_gross_amount DECIMAL(15, 2),
    
    qr_code_url TEXT,
    qr_code_string TEXT,
    qr_code_expires_at TIMESTAMP,
    qr_code_acquirer VARCHAR(50),
    
    va_numbers JSONB,
    
    notification_received_at TIMESTAMP,
    notification_payload JSONB,
    last_notification_at TIMESTAMP,
    
    notes TEXT,
    metadata JSONB,
    served_by INTEGER REFERENCES users(id),
    
    paid_at TIMESTAMP,
    expired_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_organization ON transactions(organization_id);
CREATE INDEX IF NOT EXISTS idx_transactions_outlet ON transactions(outlet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(payment_status);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_number ON transactions(transaction_number);
CREATE INDEX IF NOT EXISTS idx_transactions_midtrans_id ON transactions(midtrans_transaction_id);
CREATE INDEX IF NOT EXISTS idx_transactions_midtrans_order ON transactions(midtrans_order_id);
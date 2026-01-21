CREATE TABLE IF NOT EXISTS payment_gateways (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    outlet_id INTEGER REFERENCES outlets(id) ON DELETE CASCADE,
    environment VARCHAR(20) DEFAULT 'sandbox',
    server_key_encrypted TEXT NOT NULL,
    client_key_encrypted TEXT,
    merchant_id_midtrans VARCHAR(255),
    notification_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_gateways_organization ON payment_gateways(organization_id);
CREATE INDEX IF NOT EXISTS idx_payment_gateways_outlet ON payment_gateways(outlet_id);
CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    business_name VARCHAR(255) NOT NULL,
    business_type VARCHAR(100),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    address TEXT,
    province VARCHAR(100),
    city VARCHAR(100),
    regency VARCHAR(100),
    village VARCHAR(100),
    rt VARCHAR(5),
    rw VARCHAR(5),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'ID',
    npwp VARCHAR(50),
    subscription_status VARCHAR(50) NOT NULL DEFAULT 'trial',
    subscription_started_at TIMESTAMP,
    subscription_ends_at TIMESTAMP,
    total_outlets INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_organizations_email
ON organizations (email);

CREATE INDEX idx_organizations_subscription
ON organizations (subscription_status, subscription_ends_at);

CREATE INDEX idx_organizations_created_at
ON organizations (created_at);

CREATE INDEX idx_organizations_active
ON organizations (is_active);

ALTER TABLE organizations
ADD CONSTRAINT chk_subscription_status
CHECK (subscription_status IN (
    'trial',
    'active',
    'past_due',
    'suspended',
    'cancelled'
));

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_organizations_updated_at
BEFORE UPDATE ON organizations
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
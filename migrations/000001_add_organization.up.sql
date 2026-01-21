CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    business_name VARCHAR(255) NOT NULL,
    business_type VARCHAR(100),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'ID',
    tax_id VARCHAR(50),
    subscription_status VARCHAR(50) DEFAULT 'trial',
    subscription_started_at TIMESTAMP,
    subscription_ends_at TIMESTAMP,
    total_outlets INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_organizations_email ON organizations(email);
CREATE INDEX IF NOT EXISTS idx_organizations_subscription ON organizations(subscription_status, subscription_ends_at);

ALTER TABLE users
ADD COLUMN organization_id INT REFERENCES organization(id);
CREATE INDEX idx_users_organization_id ON users(organization_id);

UPDATE users
SET organization_id = NULL; 
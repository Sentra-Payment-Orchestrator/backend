CREATE TABLE organization (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    tax_id VARCHAR(50) UNIQUE,
    phone_number VARCHAR(15),
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_organization_name ON organization(name);
CREATE INDEX idx_organization_tax_id ON organization(tax_id);
CREATE INDEX idx_organization_email ON organization(email);

ALTER TABLE users
ADD COLUMN organization_id INT REFERENCES organization(id);
CREATE INDEX idx_users_organization_id ON users(organization_id);

UPDATE users
SET organization_id = NULL; 
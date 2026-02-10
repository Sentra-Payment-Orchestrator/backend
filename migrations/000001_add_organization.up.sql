CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    org_name VARCHAR(255) NOT NULL,
    org_type VARCHAR(100),
    org_email VARCHAR(255) NOT NULL,
    org_phone VARCHAR(50),
    org_tax_id VARCHAR(50),
    org_status SMALLINT DEFAULT 0,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_organizations_email ON organizations(org_email);

CREATE TABLE IF NOT EXISTS organization_addresses (
    id SERIAL PRIMARY KEY,
    org_id INT REFERENCES organizations(id) ON DELETE CASCADE,
    org_address TEXT NOT NULL,
    org_city VARCHAR(100) NOT NULL,
    org_state VARCHAR(100),
    org_zip_code VARCHAR(20),
    org_country VARCHAR(2) DEFAULT 'ID',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_organization_addresses_organization_id ON organization_addresses(org_id);

CREATE TABLE IF NOT EXISTS organization_roles (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL,
    role_description TEXT,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_organization_roles_id ON organization_roles(id);

INSERT INTO organization_roles (role_name, role_description, created_at, updated_at)
VALUES ('org_owner', 'Organization Owner with full access', NOW(), NOW()),
       ('org_admin', 'Organization Administrator with elevated access', NOW(), NOW()),
       ('org_manager', 'Organization Manager with management access', NOW(), NOW()),
       ('org_accountant', 'Organization Accountant with financial access', NOW(), NOW()),
       ('org_staff', 'Organization Staff with limited access', NOW(), NOW()); 

-- Organization-level user roles (many-to-many: user can be in multiple orgs with different roles)
CREATE TABLE IF NOT EXISTS organization_users (
    id SERIAL PRIMARY KEY,
    org_id INT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role INT REFERENCES organization_roles(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(org_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_organization_users_org_id ON organization_users(org_id);
CREATE INDEX IF NOT EXISTS idx_organization_users_user_id ON organization_users(user_id);
CREATE INDEX IF NOT EXISTS idx_organization_users_role ON organization_users(role);
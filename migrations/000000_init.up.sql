-- User
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,			
    email VARCHAR(255) UNIQUE NOT NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    password VARCHAR(255) NOT NULL,
    status SMALLINT DEFAULT 0 NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_username ON users(username);

CREATE TABLE user_profile (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    first_name VARCHAR(20),
    last_name VARCHAR(20),
    phone_number VARCHAR(15),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_profile_user_id ON user_profile(user_id);
CREATE INDEX idx_user_profile_email ON user_profile(email);

CREATE TABLE user_address (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    street_address VARCHAR(100),
    city VARCHAR(20),
    province VARCHAR(20),
    postal_code VARCHAR(10),
    country_code VARCHAR(2),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_address_user_id ON user_address(user_id);

CREATE TABLE role (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_role_role_name ON role(role_name);
CREATE INDEX idx_role_id ON role(id);

CREATE TABLE user_role (
    user_id INT REFERENCES users(id),
    role_id INT REFERENCES role(id),
    assigned_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);
CREATE INDEX idx_user_role_user_id ON user_role(user_id);
CREATE INDEX idx_user_role_role_id ON user_role(role_id);
CREATE INDEX idx_user_role_composite ON user_role(user_id, role_id);

CREATE TABLE permission (
    id SERIAL PRIMARY KEY,
    permission_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_permission_permission_name ON permission(permission_name);
CREATE INDEX idx_permission_id ON permission(id);

CREATE TABLE role_permission (
    role_id INT REFERENCES role(id),
    permission_id INT REFERENCES permission(id),
    assigned_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);
CREATE INDEX idx_role_permission_role_id ON role_permission(role_id);
CREATE INDEX idx_role_permission_permission_id ON role_permission(permission_id);
CREATE INDEX idx_role_permission_composite ON role_permission(role_id, permission_id);

CREATE VIEW vw_user_permissions AS
SELECT 
    u.id AS user_id,
    r.role_name,
    ARRAY_AGG(p.permission_name) AS permissions
FROM users u
JOIN user_role ur ON u.id = ur.user_id
JOIN role r ON ur.role_id = r.id
JOIN role_permission rp ON r.id = rp.role_id
JOIN permission p ON rp.permission_id = p.id
GROUP BY u.id, r.role_name;
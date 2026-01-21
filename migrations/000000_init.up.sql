-- User
CREATE TABLE users (
    id SERIAL PRIMARY KEY,			
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE user_profile (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    first_name VARCHAR(20),
    last_name VARCHAR(20),
    phone_number VARCHAR(15),
    npwp VARCHAR(20),
    birth_date DATE,
    gender VARCHAR(10),
    id_number VARCHAR(50),
    id_type VARCHAR(20),
    citizenship VARCHAR(2),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP

    CONSTRAINT fk_user_profile_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
CREATE UNIQUE INDEX ux_user_profile_user
ON user_profile(user_id);

CREATE TABLE user_address (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    address VARCHAR(100),
    city VARCHAR(20),
    province VARCHAR(20),
    regency VARCHAR(20),
    village VARCHAR(20),
    rt VARCHAR(5),
    rw VARCHAR(5),
    postal_code VARCHAR(10),
    country_code VARCHAR(2),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP

    CONSTRAINT fk_user_address_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
CREATE INDEX idx_user_address_user_id
ON user_address(user_id);
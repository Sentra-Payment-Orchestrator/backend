CREATE TABLE user_organizations (
    user_id INT NOT NULL,
    organization_id INT NOT NULL,

    role VARCHAR(50) NOT NULL DEFAULT 'member',
    status VARCHAR(50) NOT NULL DEFAULT 'active',

    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, organization_id),

    CONSTRAINT fk_uo_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_uo_organization
        FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_uo_user_id
ON user_organizations(user_id);

CREATE INDEX idx_uo_organization_id
ON user_organizations(organization_id);

CREATE INDEX idx_uo_role
ON user_organizations(role);
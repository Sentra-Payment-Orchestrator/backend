CREATE TABLE IF NOT EXISTS outlets (
    id SERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    nitku VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    phone VARCHAR(50),
    email VARCHAR(255),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_outlets_org_id ON outlets(org_id);
CREATE INDEX IF NOT EXISTS idx_outlets_id ON outlets(id);

-- Outlet-level user roles (many-to-many: user can work at multiple outlets with different roles)
-- outlet_role examples: outlet_manager, outlet_supervisor, outlet_cashier, outlet_staff
CREATE TABLE IF NOT EXISTS outlet_users (
    id SERIAL PRIMARY KEY,
    outlet_id INTEGER NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    outlet_role VARCHAR(50) DEFAULT 'outlet_staff', -- outlet_manager, outlet_supervisor, outlet_cashier, outlet_staff
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(outlet_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_outlet_users_outlet_id ON outlet_users(outlet_id);
CREATE INDEX IF NOT EXISTS idx_outlet_users_user_id ON outlet_users(user_id);
CREATE INDEX IF NOT EXISTS idx_outlet_users_role ON outlet_users(outlet_role);
CREATE INDEX IF NOT EXISTS idx_outlet_users_user_id ON outlet_users(user_id);

CREATE TABLE IF NOT EXISTS outlet_overview (
    id SERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    outlet_id INTEGER NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    subscription_start_date TIMESTAMPTZ,
    subscription_end_date TIMESTAMPTZ,
    assigned_user_count INTEGER DEFAULT 0,
    total_transactions INTEGER DEFAULT 0,
    last_transaction_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_outlet_overview_org_id ON outlet_overview(org_id);
CREATE INDEX IF NOT EXISTS idx_outlet_overview_outlet ON outlet_overview(outlet_id);

CREATE OR REPLACE FUNCTION update_outlet_overview_on_insert()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO outlet_overview (org_id, outlet_id, created_at, updated_at)
    VALUES (NEW.org_id, NEW.id, NOW(), NOW());

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_outlet_overview_insert
AFTER INSERT ON outlets
FOR EACH ROW
EXECUTE FUNCTION update_outlet_overview_on_insert();

-- Trigger to update outlet_overview when outlet_users changes
CREATE OR REPLACE FUNCTION update_outlet_user_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE outlet_overview
        SET assigned_user_count = assigned_user_count + 1,
            updated_at = NOW()
        WHERE outlet_id = NEW.outlet_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE outlet_overview
        SET assigned_user_count = GREATEST(assigned_user_count - 1, 0),
            updated_at = NOW()
        WHERE outlet_id = OLD.outlet_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_outlet_user_count
AFTER INSERT OR DELETE ON outlet_users
FOR EACH ROW
EXECUTE FUNCTION update_outlet_user_count();
CREATE TABLE IF NOT EXISTS daily_sales_summary (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    outlet_id INTEGER REFERENCES outlets(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_transactions INTEGER DEFAULT 0,
    total_revenue DECIMAL(15, 2) DEFAULT 0,
    total_cost DECIMAL(15, 2) DEFAULT 0,
    total_profit DECIMAL(15, 2) DEFAULT 0,
    top_selling_products JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(organization_id, outlet_id, date)
);

CREATE INDEX IF NOT EXISTS idx_daily_sales_organization ON daily_sales_summary(organization_id, date);
CREATE INDEX IF NOT EXISTS idx_daily_sales_outlet ON daily_sales_summary(outlet_id, date);
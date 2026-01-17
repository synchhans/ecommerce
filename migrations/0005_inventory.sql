-- ===== Inventory (single warehouse) =====
CREATE TABLE IF NOT EXISTS inventory_items (
  variant_id uuid PRIMARY KEY REFERENCES product_variants(id) ON DELETE CASCADE,
  stock_on_hand int NOT NULL DEFAULT 0 CHECK (stock_on_hand >= 0),
  reserved int NOT NULL DEFAULT 0 CHECK (reserved >= 0),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_inventory_updated ON inventory_items(updated_at DESC);

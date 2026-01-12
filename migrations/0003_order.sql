-- ===== Orders =====
CREATE TABLE IF NOT EXISTS orders (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_number text UNIQUE NOT NULL,
  user_id uuid NULL,
  cart_id uuid NULL REFERENCES carts(id) ON DELETE SET NULL,
  status text NOT NULL DEFAULT 'pending_payment',
  currency text NOT NULL DEFAULT 'IDR',

  subtotal bigint NOT NULL CHECK (subtotal >= 0),
  discount_total bigint NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
  shipping_total bigint NOT NULL DEFAULT 0 CHECK (shipping_total >= 0),
  grand_total bigint NOT NULL CHECK (grand_total >= 0),

  shipping_address_snapshot jsonb NOT NULL DEFAULT '{}'::jsonb,

  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS order_items (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  variant_id uuid NOT NULL REFERENCES product_variants(id) ON DELETE RESTRICT,

  sku text NOT NULL,
  name text NOT NULL,
  unit_price bigint NOT NULL CHECK (unit_price >= 0),
  qty int NOT NULL CHECK (qty > 0),
  line_total bigint NOT NULL CHECK (line_total >= 0)
);

CREATE INDEX IF NOT EXISTS idx_orders_created ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);

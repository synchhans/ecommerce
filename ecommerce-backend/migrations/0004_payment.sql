-- ===== Payments =====
CREATE TABLE IF NOT EXISTS payments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

  provider text NOT NULL, -- e.g. "manual", "midtrans", "xendit"
  status text NOT NULL DEFAULT 'initiated', -- initiated/pending/paid/failed/expired/refunded
  amount bigint NOT NULL CHECK (amount >= 0),

  provider_ref text NULL,  -- reference id from provider (or internal ref)
  pay_url text NULL,       -- optional redirect url
  payload jsonb NULL,      -- raw callback payload

  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payments_order ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_provider_ref ON payments(provider, provider_ref);

-- ===== Idempotency Keys (optional but recommended) =====
CREATE TABLE IF NOT EXISTS idempotency_keys (
  key text PRIMARY KEY,
  request_hash text NOT NULL,
  response_body jsonb NOT NULL,
  status_code int NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

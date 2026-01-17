-- ===== User Addresses =====
CREATE TABLE IF NOT EXISTS user_addresses (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  label text NOT NULL DEFAULT '',
  recipient_name text NOT NULL,
  phone text NOT NULL,

  address_line1 text NOT NULL,
  address_line2 text NULL,

  city text NOT NULL,
  province text NOT NULL,
  postal_code text NOT NULL,
  country text NOT NULL DEFAULT 'ID',

  is_default boolean NOT NULL DEFAULT false,

  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_user_addresses_user ON user_addresses(user_id);

-- Only one default address per user
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_addresses_default_per_user
ON user_addresses(user_id)
WHERE is_default = true;

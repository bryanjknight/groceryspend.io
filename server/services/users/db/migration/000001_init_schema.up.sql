CREATE TABLE IF NOT EXISTS "user" (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  auth_provider_id VARCHAR(255) NOT NULL,
  auth_specific_id VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "auth_provider_idx" on "user" (auth_provider_id, auth_specific_id);
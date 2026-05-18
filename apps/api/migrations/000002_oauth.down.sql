DROP INDEX IF EXISTS idx_users_provider_provider_id;

ALTER TABLE users
    DROP COLUMN IF EXISTS provider_id,
    DROP COLUMN IF EXISTS provider,
    ALTER COLUMN password_hash SET NOT NULL;

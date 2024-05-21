DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'genders') THEN
    CREATE TYPE genders AS ENUM ('male', 'female');
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS patients (
  id VARCHAR(26) NOT NULL PRIMARY KEY,
  identity_number VARCHAR(16) NOT NULL UNIQUE,
  phone_number VARCHAR(15) NOT NULL,
  name VARCHAR(30) NOT NULL,
  birth_date TIMESTAMP NOT NULL,
  gender genders NOT NULL,
  card_image_url TEXT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "line_id" varchar NOT NULL,
  "nickname" varchar NOT NULL,
  "is_email_verified" bool NOT NULL DEFAULT false,
  "role" varchar NOT NULL DEFAULT 'unverified'
);

CREATE INDEX ON "users" ("email");

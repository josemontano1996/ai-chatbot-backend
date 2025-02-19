CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT 'gen_random_uuid()',
  "email" varchar(100) UNIQUE NOT NULL,
  "password" varchar(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

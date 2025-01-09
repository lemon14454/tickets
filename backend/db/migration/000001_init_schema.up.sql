CREATE TYPE "event_status" AS ENUM (
  'processing',
  'created',
  'available',
  'done',
  'failure'
);

CREATE TABLE "events" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "host_id" bigint NOT NULL,
  "start_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "status" event_status NOT NULL DEFAULT 'processing'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(20) NOT NULL,
  "email" varchar(50) NOT NULL,
  "hashed_password" varchar(255) NOT NULL,
  "host" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tickets" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "order_id" bigint,
  "event_id" bigint NOT NULL,
  "zone_id" bigint NOT NULL,
  "row" int NOT NULL,
  "seat" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  UNIQUE(event_id, zone_id, row, seat)
);

CREATE TABLE "event_zones" (
  "id" bigserial PRIMARY KEY,
  "zone" varchar(10) NOT NULL,
  "event_id" bigint NOT NULL,
  "rows" int NOT NULL,
  "seats" int NOT NULL,
  "price" int NOT NULL DEFAULT 1000
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "event_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "total_price" int NOT NULL
);

ALTER TABLE "events" ADD FOREIGN KEY ("host_id") REFERENCES "users" ("id");
ALTER TABLE "tickets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "tickets" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
ALTER TABLE "tickets" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");
ALTER TABLE "tickets" ADD FOREIGN KEY ("zone_id") REFERENCES "event_zones" ("id");
ALTER TABLE "event_zones" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");
ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "orders" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");

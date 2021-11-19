CREATE TABLE "User" (
  "id" SERIAL PRIMARY KEY,
  "email" string UNIQUE,
  "firstname" varchar,
  "lastname" varchar,
  "password" varchar,
  "merchant_id" int,
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "Merchant" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "logo" varchar
);

CREATE TABLE "Outlet" (
  "id" SERIAL PRIMARY KEY,
  "merchant_id" int,
  "name" varchar,
  "location" varchar
);

CREATE TABLE "Product" (
  "id" SERIAL PRIMARY KEY,
  "merchant_id" int,
  "sku" varchar,
  "name" varchar,
  "display_image" varchar
);

CREATE TABLE "OutletProduct" (
  "outlet_id" int,
  "product_id" int,
  "price" float4,
  "stock" int
);

ALTER TABLE "User" ADD FOREIGN KEY ("merchant_id") REFERENCES "Merchant" ("id");

ALTER TABLE "Outlet" ADD FOREIGN KEY ("merchant_id") REFERENCES "Merchant" ("id");

ALTER TABLE "Product" ADD FOREIGN KEY ("merchant_id") REFERENCES "Merchant" ("id");

ALTER TABLE "OutletProduct" ADD FOREIGN KEY ("outlet_id") REFERENCES "Outlet" ("id");

ALTER TABLE "OutletProduct" ADD FOREIGN KEY ("product_id") REFERENCES "Product" ("id");

CREATE INDEX ON "User" ("id");

CREATE INDEX ON "User" ("email");

CREATE INDEX ON "User" ("merchant_id");

CREATE INDEX ON "Merchant" ("id");

CREATE INDEX ON "Merchant" ("name");

CREATE INDEX ON "Outlet" ("id");

CREATE INDEX ON "Outlet" ("merchant_id");

CREATE INDEX ON "Outlet" ("name");

CREATE INDEX ON "Product" ("id");

CREATE INDEX ON "Product" ("merchant_id");

CREATE INDEX ON "Product" ("sku");

CREATE INDEX ON "Product" ("name");

CREATE INDEX ON "OutletProduct" ("outlet_id");

CREATE INDEX ON "OutletProduct" ("product_id");

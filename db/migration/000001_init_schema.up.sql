CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "name" varchar NOT NULL,
                         "email" varchar UNIQUE NOT NULL,
                         "username" varchar NOT NULL,
                         "bank_code" bigint DEFAULT 0,
                         "password" varchar,
                         "balance" decimal(16,2) DEFAULT (0),
                         "phone" varchar NOT NULL,
                         "identity_number" varchar NOT NULL,
                         "verified_at" timestamptz DEFAULT (null),
                         "created_at" timestamptz DEFAULT (now()),
                         "updated_at" timestamptz DEFAULT (null),
                         "deleted_at" timestamptz DEFAULT (null),
                         "created_by" bigint DEFAULT (null),
                         "updated_by" bigint DEFAULT (null),
                         "deleted_by" bigint DEFAULT (null)
);

CREATE TABLE "roles" (
                         "id" bigserial PRIMARY KEY,
                         "name" varchar NOT NULL,
                         "level" smallint NOT NULL,
                         "created_at" timestamptz DEFAULT (now()),
                         "updated_at" timestamptz DEFAULT (null),
                         "deleted_at" timestamptz DEFAULT (null),
                         "created_by" bigint DEFAULT (null),
                         "updated_by" bigint DEFAULT (null),
                         "deleted_by" bigint DEFAULT (null)
);

CREATE TABLE "role_users" (
                              "id" bigserial PRIMARY KEY,
                              "role_id" bigint NOT NULL,
                              "user_id" bigint NOT NULL,
                              "created_at" timestamptz DEFAULT (now()),
                              "updated_at" timestamptz DEFAULT (null),
                              "deleted_at" timestamptz DEFAULT (null),
                              "created_by" bigint DEFAULT (null),
                              "updated_by" bigint DEFAULT (null),
                              "deleted_by" bigint DEFAULT (null)
);

CREATE TABLE "categories" (
                              "id" bigserial PRIMARY KEY,
                              "name" varchar(50) NOT NULL,
                              "up_selling" decimal(16,2) DEFAULT 0,
                              "parent" bigint NOT NULL DEFAULT (0),
                              "created_at" timestamptz DEFAULT (now()),
                              "updated_at" timestamptz DEFAULT (null),
                              "deleted_at" timestamptz DEFAULT (null),
                              "created_by" bigint DEFAULT (null),
                              "updated_by" bigint DEFAULT (null),
                              "deleted_by" bigint DEFAULT (null)
);

CREATE TABLE "sellings" (
                            "id" bigserial PRIMARY KEY,
                            "partner_id" bigint DEFAULT 0,
                            "category_id" bigint DEFAULT 0,
                            "amount" decimal(16,2) DEFAULT 0,
                            "created_at" timestamptz DEFAULT (now()),
                            "updated_at" timestamptz DEFAULT (null),
                            "deleted_at" timestamptz DEFAULT (null),
                            "created_by" bigint DEFAULT 0,
                            "updated_by" bigint DEFAULT (null),
                            "deleted_by" bigint DEFAULT (null)
);

CREATE TABLE "products" (
                            "id" bigserial PRIMARY KEY,
                            "cat_id" bigint NOT NULL,
                            "name" varchar(50) NOT NULL,
                            "amount" decimal(16,2) NOT NULL DEFAULT 0,
                            "provider_id" bigint NOT NULL,
                            "provider_code" varchar(50) NOT NULL,
                            "status" varchar(20) NOT NULL,
                            "parent" bigint NOT NULL DEFAULT (0),
                            "created_at" timestamptz DEFAULT (now()),
                            "updated_at" timestamptz DEFAULT (null),
                            "deleted_at" timestamptz DEFAULT (null),
                            "created_by" bigint DEFAULT (null),
                            "updated_by" bigint DEFAULT (null),
                            "deleted_by" bigint DEFAULT (null)
);

CREATE TABLE "providers" (
                             "id" bigserial PRIMARY KEY,
                             "name" varchar(50) NOT NULL,
                             "user" varchar(100) NOT NULL,
                             "secret" varchar(100) NOT NULL,
                             "add_info1" varchar(150) NOT NULL,
                             "add_info2" varchar(150) NOT NULL,
                             "valid_from" timestamptz DEFAULT (now()),
                             "valid_to" timestamptz DEFAULT (null),
                             "base_url" varchar(100),
                             "method" varchar(10),
                             "inq" varchar(50),
                             "pay" varchar(50),
                             "adv" varchar(50),
                             "cmt" varchar(50),
                             "rev" varchar(50),
                             "status" varchar(20) NOT NULL,
                             "created_at" timestamptz DEFAULT (now()),
                             "updated_at" timestamptz DEFAULT (null),
                             "deleted_at" timestamptz DEFAULT (null),
                             "created_by" bigint DEFAULT (null),
                             "updated_by" bigint DEFAULT (null),
                             "deleted_by" bigint DEFAULT (null)
);

CREATE TABLE "partners" (
                            "id" bigserial PRIMARY KEY,
                            "name" varchar(50) NOT NULL,
                            "user" varchar(100) NOT NULL,
                            "secret" varchar(100) NOT NULL,
                            "add_info1" varchar(150) NOT NULL,
                            "add_info2" varchar(150) NOT NULL,
                            "valid_from" timestamptz DEFAULT (now()),
                            "valid_to" timestamptz DEFAULT (null),
                            "payment_type" varchar(20) NOT NULL,
                            "status" varchar(20) NOT NULL,
                            "created_at" timestamptz DEFAULT (now()),
                            "updated_at" timestamptz DEFAULT (null),
                            "deleted_at" timestamptz DEFAULT (null),
                            "created_by" bigint DEFAULT (null),
                            "updated_by" bigint DEFAULT (null),
                            "deleted_by" bigint DEFAULT (null)
);

CREATE TABLE "transactions" (
                                "id" bigserial PRIMARY KEY,
                                "tx_id" varchar(50) NOT NULL,
                                "ref_id" varchar(50) NOT NULL,
                                "bill_id" varchar(50) NOT NULL,
                                "cust_name" varchar DEFAULT null,
                                "amount" decimal(16,2) DEFAULT 0,
                                "admin" decimal(16,2) DEFAULT 0,
                                "tot_amount" decimal(16,2) DEFAULT 0,
                                "fee_partner" decimal(16,2) DEFAULT 0,
                                "fee_ppob" decimal(16,2) DEFAULT 0,
                                "first_balance" decimal(16,2) DEFAULT 0,
                                "last_balance" decimal(16,2) DEFAULT 0,
                                "valid_from" timestamptz DEFAULT (now()),
                                "valid_to" timestamptz DEFAULT (null),
                                "cat_id" bigint DEFAULT null,
                                "cat_name" varchar(100) DEFAULT null,
                                "prod_id" bigint DEFAULT null,
                                "prod_name" varchar(100) DEFAULT null,
                                "partner_id" bigint DEFAULT null,
                                "partner_name" varchar(100) DEFAULT null,
                                "provider_id" bigint DEFAULT null,
                                "provider_name" varchar(100) DEFAULT null,
                                "status" varchar(20) NOT NULL,
                                "req_inq_params" text DEFAULT null,
                                "res_inq_params" text DEFAULT null,
                                "req_pay_params" text DEFAULT null,
                                "res_pay_params" text DEFAULT null,
                                "req_cmt_params" text DEFAULT null,
                                "res_cmt_params" text DEFAULT null,
                                "req_adv_params" text DEFAULT null,
                                "res_adv_params" text DEFAULT null,
                                "req_rev_params" text DEFAULT null,
                                "res_rev_params" text DEFAULT null,
                                "created_at" timestamptz DEFAULT (now()),
                                "updated_at" timestamptz DEFAULT (null),
                                "deleted_at" timestamptz DEFAULT (null),
                                "created_by" bigint DEFAULT (null),
                                "updated_by" bigint DEFAULT (null),
                                "deleted_by" bigint DEFAULT (null)
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email", "username");

CREATE INDEX ON "roles" ("name");

CREATE INDEX ON "roles" ("level");

CREATE UNIQUE INDEX ON "roles" ("name", "level");

CREATE INDEX ON "role_users" ("role_id");

CREATE INDEX ON "role_users" ("user_id");

CREATE UNIQUE INDEX ON "role_users" ("role_id", "user_id");

CREATE INDEX ON "categories" ("name");

CREATE INDEX ON "categories" ("parent");

CREATE UNIQUE INDEX ON "categories" ("name", "parent");

CREATE INDEX ON "sellings" ("partner_id");

CREATE INDEX ON "sellings" ("category_id");

CREATE UNIQUE INDEX ON "sellings" ("partner_id", "category_id");

CREATE INDEX ON "products" ("cat_id");

CREATE INDEX ON "products" ("provider_id");

CREATE INDEX ON "products" ("cat_id", "provider_id");

CREATE INDEX ON "providers" ("name");

CREATE INDEX ON "partners" ("name");

CREATE INDEX ON "transactions" ("cat_id");

CREATE INDEX ON "transactions" ("tx_id");

CREATE INDEX ON "transactions" ("ref_id");

CREATE INDEX ON "transactions" ("prod_id");

CREATE INDEX ON "transactions" ("provider_id");

CREATE INDEX ON "transactions" ("status");

CREATE INDEX ON "transactions" ("created_at");

CREATE INDEX ON "transactions" ("status", "prod_id");

CREATE INDEX ON "transactions" ("status", "provider_id");

CREATE INDEX ON "transactions" ("cat_id", "prod_id");

CREATE INDEX ON "transactions" ("created_at", "status");

ALTER TABLE "users" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "roles" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "roles" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "roles" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "role_users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "role_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "role_users" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "role_users" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "role_users" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "categories" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "categories" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "categories" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "sellings" ADD FOREIGN KEY ("partner_id") REFERENCES "products" ("id");

ALTER TABLE "sellings" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "sellings" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "sellings" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "sellings" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("cat_id") REFERENCES "categories" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("provider_id") REFERENCES "providers" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "providers" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "providers" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "providers" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "partners" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "partners" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "partners" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("cat_id") REFERENCES "categories" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("prod_id") REFERENCES "products" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("partner_id") REFERENCES "partners" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("provider_id") REFERENCES "providers" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

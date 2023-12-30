CREATE TABLE "mentor" (
    "id" bigserial PRIMARY KEY,
    "pic" varchar(2048),
    "username" varchar UNIQUE NOT NULL,
    "hashed_password" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password_change_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "career" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "sector" varchar NOT NULL,
    "eta" varchar NOT NULL,
    "blog" varchar(2048),
    "estimated_income" int NOT NULL
);


ALTER TABLE "career" ADD FOREIGN KEY ("id") REFERENCES "mentor" ("id");
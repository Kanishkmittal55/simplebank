CREATE TABLE "otp_emails" (
      "id" bigserial PRIMARY KEY,
      "username" varchar NOT NULL,
      "email" varchar NOT NULL,
      "otp_sent" varchar NOT NULL,
      "total_otp_generations" int NOT NULL DEFAULT 0,
      "created_at" timestamptz NOT NULL DEFAULT (now()),
      "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE "otp_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
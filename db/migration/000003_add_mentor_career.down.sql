-- First, remove the foreign key constraints
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS accounts_owner_idx;

ALTER TABLE "career" DROP CONSTRAINT IF EXISTS career_id_fkey;

-- Now, drop the tables
DROP TABLE IF EXISTS "career" CASCADE;
DROP TABLE IF EXISTS "mentor" CASCADE;

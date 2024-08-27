CREATE TABLE  IF NOT EXISTS  "pets" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "gender" varchar,
  "user_id" integer,
  "dob" varchar,
  "neutered" bool,
  "vaccinated" bool,
  "created_at" timestamp
);
ALTER TABLE "users" ADD "petId" integer;
ALTER TABLE "users" ADD FOREIGN KEY ("petId") REFERENCES "pets" ("id");
ALTER TABLE "pets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

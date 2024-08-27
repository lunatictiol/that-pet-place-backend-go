ALTER TABLE "pets" DROP FOREIGN KEY "pets_user_id_fkey";
ALTER TABLE "users" DROP FOREIGN KEY "users_petid_fkey";
DROP TABLE "pets";
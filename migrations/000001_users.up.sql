CREATE TABLE "Users" (
  Id INTEGER PRIMARY KEY,
  Login VARCHAR(255) UNIQUE NOT NULL,
  PasswordHash VARCHAR(255) NOT NULL
)
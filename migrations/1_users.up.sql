CREATE TABLE users(
   id UUID PRIMARY KEY,
   email TEXT UNIQUE NOT NULL,
   password TEXT NOT NULL,
   created_at TIMESTAMP NOT NULL,
   last_login TIMESTAMP
);
CREATE TABLE IF NOT EXISTS "user" (
	id serial PRIMARY KEY,
	username varchar(30) UNIQUE NOT NULL,
	email varchar(30) UNIQUE NOT NULL,
	password varchar(60) NOT NULL,
	created_at timestamp(2) NOT NULL default CURRENT_TIMESTAMP (2) 
);

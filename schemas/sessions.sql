CREATE TABLE IF NOT EXISTS "sessions" (
	uuid varchar(36) PRIMARY KEY,
	email varchar(30) UNIQUE NOT NULL,
	expires_at timestamp(2) NOT NULL default CURRENT_TIMESTAMP (2) 
);

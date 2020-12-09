package crdb

var schema = `
USE defaultdb;
CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY UNIQUE,
	username STRING UNIQUE,
	password STRING,
	email STRING,
	bio STRING,
	credit INT,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);
`

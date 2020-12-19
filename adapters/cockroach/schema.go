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

CREATE TABLE IF NOT EXISTS contents(
	id UUID PRIMARY KEY,
	cid STRING UNIQUE,
	uploader_id UUID REFERENCES users(id),
	name STRING,
	description STRING,
	filename STRING,
	extension STRING,
	category STRING,
	size FLOAT,
	downloads INT
);
`

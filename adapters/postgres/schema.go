package postgres

// TODO: remove the first lines
var schema = `
DROP TABLE ratings;
DROP TABLE contents; 
DROP TABLE users;
DROP TABLE ftypes;

CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY,
	username varchar(40) NOT NULL UNIQUE,
	password char(60) NOT NULL,
	email text NOT NULL,
	credit INT NOT NULL DEFAULT 0,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ftypes(
	id serial PRIMARY KEY,
	file_type varchar(15) NOT NULL
);

INSERT INTO ftypes(file_type) VALUES ('font'),('text'),('image'),('audio'),('video'),
('spreadsheet'),('presentation'),('document'),('archive'),('application');

CREATE TABLE IF NOT EXISTS contents(
	id UUID PRIMARY KEY,
	cid text UNIQUE NOT NULL,
	uploader_id UUID REFERENCES users(id) ON DELETE RESTRICT,
	name varchar(75) NOT NULL,
	extension varchar(10) NOT NULL,
	type_id INT NOT NULL REFERENCES ftypes(id) ON DELETE RESTRICT,
	description varchar(200),
	size FLOAT NOT NULL,
	downloads INT NOT NULL DEFAULT 0,
	rating FLOAT check(rating >= 0 and rating <= 5) NOT NULL DEFAULT 0,
	uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ratings(
	rating FLOAT check(rating >= 0 and rating <= 5) NOT NULL,
	user_id UUID REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
	content_id UUID REFERENCES contents(id) ON DELETE RESTRICT NOT NULL,
	CONSTRAINT unique_ratings UNIQUE(content_id, user_id)
);

CREATE OR REPLACE Function update_rating() RETURNS trigger AS $update_rating$
BEGIN
update contents set rating = (select avg(rating) from ratings where content_id=NEW.content_id)
where id=NEW.content_id;
return null;
END;
$update_rating$ LANGUAGE plpgsql;

CREATE TRIGGER update_rating AFTER INSERT OR UPDATE ON ratings
FOR EACH ROW EXECUTE FUNCTION update_rating();

`

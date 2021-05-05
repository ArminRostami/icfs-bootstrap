

-- DROP TABLE downloads;
-- DROP TABLE contents; 
-- DROP TABLE users;
-- DROP TABLE ftypes;



CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY,
	username varchar(40) NOT NULL UNIQUE,
	password char(60) NOT NULL,
	email text NOT NULL UNIQUE,
	credit INT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ftypes(
	id serial PRIMARY KEY,
	file_type varchar(15) UNIQUE NOT NULL
);

CREATE OR REPLACE FUNCTION create_types() RETURNS void AS $$
declare 
	type_count int;
begin
	select count(*) into type_count from ftypes;
	if type_count < 10 THEN
	INSERT INTO ftypes(file_type) VALUES ('font'),('text'),('image'),('audio'),('video'),
	('spreadsheet'),('presentation'),('document'),('archive'),('application');
	raise notice 'file types created.';
	end if;
end
$$ LANGUAGE plpgsql;

SELECT create_types();


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
	rating FLOAT check(rating >= 0 and rating <= 5) NOT NULL DEFAULT 2.5,
	uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_modified TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

	tsv tsvector GENERATED ALWAYS AS (
		setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
		setweight(to_tsvector('english', coalesce(extension, '')), 'B') ||
		setweight(to_tsvector('english', coalesce(description, '')), 'B') 
	) STORED
);



CREATE INDEX IF NOT EXISTS textsearch_idx ON contents USING GIN (tsv);

CREATE TABLE IF NOT EXISTS downloads(
	user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
	content_id UUID REFERENCES contents(id) ON DELETE CASCADE NOT NULL,
	rating FLOAT check(rating >= 0 and rating <= 5) NOT NULL DEFAULT 2.5,
	comment_text varchar(200),
	comment_time TIMESTAMPTZ,
	downloaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT unique_ratings UNIQUE(content_id, user_id)
);

CREATE OR REPLACE Function update_rating() RETURNS trigger AS $update_rating$
BEGIN
update contents set rating = (select avg(rating) from downloads where content_id=NEW.content_id)
where id=NEW.content_id;
return null;
END;
$update_rating$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_rating on downloads;

CREATE TRIGGER update_rating AFTER INSERT OR UPDATE ON downloads
FOR EACH ROW EXECUTE FUNCTION update_rating();

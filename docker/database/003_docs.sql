\c doc_server
CREATE SCHEMA IF NOT EXISTS docs;

CREATE TABLE IF NOT EXISTS docs.meta
(
	id int UNIQUE NOT NULL,
	name VARCHAR(256) NOT NULL,
	file bool NOT NULL DEFAULT true,
	public bool NOT NULL DEFAULT false,
	mime   VARCHAR(256) NOT NULL,
	created_at timestamptz NOT NULL DEFAULT NOW(),
	updated_at timestamptz NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_docs_trgr BEFORE UPDATE ON docs.meta FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_docs_trgr BEFORE UPDATE ON docs.meta FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS docs.files
(
	file_id int UNIQUE NOT NULL,
	file bytea NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS docs.permissions
(
	file_id int NOT NULL,
	user_id int NOT NULL
);

GRANT USAGE ON SCHEMA docs TO doc_server_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA docs TO doc_server_admin;
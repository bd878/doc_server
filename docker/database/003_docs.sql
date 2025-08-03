\c doc_server
CREATE SCHEMA IF NOT EXISTS docs;

CREATE TABLE IF NOT EXISTS docs.meta
(
	id                 varchar(256) UNIQUE NOT NULL,
	oid                int UNIQUE DEFAULT NULL, -- large object id
	name               varchar(256) NOT NULL,
	file               bool NOT NULL DEFAULT true,
	json               bytea DEFAULT NULL,
	public             bool NOT NULL DEFAULT false,
	mime               varchar(256) NOT NULL,
	owner_login        varchar(256) NOT NULL,
	created_at         timestamptz NOT NULL DEFAULT NOW(),
	updated_at         timestamptz NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id),
	CONSTRAINT file_oid_check CHECK (file IS true AND oid IS NOT NULL)
);

CREATE TRIGGER created_at_docs_trgr BEFORE UPDATE ON docs.meta FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_docs_trgr BEFORE UPDATE ON docs.meta FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS docs.permissions
(
	file_id      varchar(256) NOT NULL REFERENCES meta(id),
	user_login   varchar(256) NOT NULL
);

GRANT USAGE ON SCHEMA docs TO doc_server_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA docs TO doc_server_admin;
CREATE DATABASE doc_server;

GRANT CONNECT ON DATABASE doc_server TO doc_server_admin;
ALTER ROLE doc_server_admin SET search_path TO doc_server, "$user", public;

\c doc_server

CREATE OR REPLACE FUNCTION created_at_trigger()
RETURNS TRIGGER AS $$
BEGIN
	NEW.created_at := OLD.created_at;
	RETURN NEW;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION updated_at_trigger()
RETURNS TRIGGER AS $$
BEGIN
	IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
		NEW.updated_at = NOW();
		RETURN NEW;
	ELSE
		RETURN OLD;
	END IF;
END;
$$ language 'plpgsql';
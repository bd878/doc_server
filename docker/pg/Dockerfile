FROM postgres:16-alpine AS database

WORKDIR /doc_server
COPY docker/database /docker-entrypoint-initdb.d/
EXPOSE 5432
CMD ["postgres"]
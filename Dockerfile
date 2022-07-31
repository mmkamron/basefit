FROM postgres
ENV POSTGRES_PASSWORD library
ENV POSTGRES_DB library
COPY table.sql /docker-entrypoint-initdb.d/

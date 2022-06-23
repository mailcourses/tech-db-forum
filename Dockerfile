FROM alpine
COPY db/db.sql /initdb/db.sql

RUN apk update
RUN apk add postgresql

RUN mkdir /run/postgresql
RUN chown postgres:postgres /run/postgresql/

RUN su postgres -
RUN mkdir /var/lib/postgresql/data
RUN chmod 0700 /var/lib/postgresql/data

RUN initdb -D /initdb/db.sql
RUN echo "host all all 0.0.0.0/0 md5" >> /var/lib/postgresql/data/pg_hba.conf
RUN echo "listen_addresses='*'" >> /var/lib/postgresql/data/postgresql.conf

RUN pg_ctl start -D /var/lib/postgresql/data
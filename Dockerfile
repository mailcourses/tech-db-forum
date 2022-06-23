FROM golang:latest AS builder

WORKDIR /app

COPY api/ /var/www/service/api/
COPY go.mod /var/www/service/go.mod
COPY go.sum /var/www/service/go.sum

WORKDIR /var/www/service/

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g api/cmd/main.go
RUN go build -ldflags "-w -s" -o ./tmp/api api/cmd/main.go

FROM ubuntu

RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo Russia/Moscow > /etc/timezone

RUN apt -y update && apt upgrade -y
RUN apt -y install gnupg2 wget vim
RUN apt-cache search postgresql | grep postgresql
RUN apt install -y postgresql postgresql-contrib
RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt -y install postgresql-14

USER postgres

RUN /etc/init.d/postgresql start && \
  psql --command "CREATE USER root WITH SUPERUSER PASSWORD 'admin';" && \
  createdb -O root forum_db && \
  /etc/init.d/postgresql stop

USER root

COPY ./db/db.sql ./db.sql

COPY --from=builder /var/www/service/tmp/api .

EXPOSE 5000
EXPOSE 5432
ENV PGPASSWORD admin
ENV dsn user=root password=$PGPASSWORD dbname=forum_db host=localhost port=5432 sslmode=disable
ENV port "5000"
ENV dbType postgres
CMD service postgresql start && psql -h localhost -d forum_db -U root -p 5432 -a -q -f ./db.sql && ./api

FROM golang:1.12-stretch AS lang

ADD . /go/src/park_base/park_db
WORKDIR /usr/src/source

COPY . .
RUN go get -d && go build -v


FROM ubuntu:18.04

ENV PGSQLVERSION 10
ENV PORT 5000

RUN apt-get update && apt-get install -y postgresql-$PGSQLVERSION

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER gd1 WITH SUPERUSER PASSWORD '123';" &&\
    createdb -O gd1 api_db &&\
    psql -d api_db -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    /etc/init.d/postgresql stop


USER root

RUN echo "listen_addresses = '*'\nsynchronous_commit = off\nfsync = off" >> /etc/postgresql/$PGSQLVERSION/main/postgresql.conf
RUN echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$PGSQLVERSION/main/postgresql.conf
   
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

EXPOSE $PORT

USER postgres

WORKDIR /usr/src/source/source

COPY --from=lang /usr/src/source/source .

COPY database/sql/ database/sql/

CMD /etc/init.d/postgresql start && ./source

FROM postgres:16-alpine

ARG dbuser
ARG dbpass
ARG dbname

ENV LANG ja_JP.utf8
ENV TZ "Asia/Tokyo"
ENV POSTGRES_USER ${dbuser}
ENV POSTGRES_PASSWORD ${dbpass}
ENV POSTGRES_DB ${dbname}

RUN mkdir ddl dml
RUN echo $POSTGRES_USER
RUN echo $POSTGRES_PASSWORD
RUN echo $POSTGRES_DB

COPY ./docker/postgresql/*.sql /docker-entrypoint-initdb.d

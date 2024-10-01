FROM flyway/flyway:10.18

# Include psql client
RUN apt update
RUN apt install postgresql-client -y

WORKDIR /migrations/

COPY . .

ENTRYPOINT [ "flyway" ]
CMD ["-configFiles=flyway.config.toml", "migrate"]

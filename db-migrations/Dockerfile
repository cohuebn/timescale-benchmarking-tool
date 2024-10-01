FROM flyway/flyway:10.18

WORKDIR /migrations/

COPY . .

ENTRYPOINT [ "flyway" ]
CMD ["-configFiles=flyway.config.toml", "migrate"]

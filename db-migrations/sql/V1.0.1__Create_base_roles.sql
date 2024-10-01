-- Create a role that can only read from the database

do $$
begin
  if not exists (select from pg_catalog.pg_roles where rolname = 'readonly') then
    create role readonly;
  end if;
end$$;

grant connect on database ${flyway:database} to readonly;
grant select on all tables in schema ${flyway:defaultSchema} to readonly;
grant select on all sequences in schema ${flyway:defaultSchema} to readonly;
grant usage on schema ${flyway:defaultSchema} to readonly;
alter default privileges in schema ${flyway:defaultSchema} grant select on tables to readonly;
alter default privileges in schema ${flyway:defaultSchema} grant select on sequences to readonly;
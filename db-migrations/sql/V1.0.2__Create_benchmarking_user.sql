do $$
begin
  if not exists (select from pg_catalog.pg_roles where rolname = '${benchmarkingUser}') then
    create role ${benchmarkingUser};
  end if;
  alter role ${benchmarkingUser} with password '${benchmarkingPassword}' login;
end$$;

grant readonly to ${benchmarkingUser}

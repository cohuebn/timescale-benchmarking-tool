create table cpu_usage(
  ts    timestamptz,
  host  text,
  usage double precision
);
select create_hypertable('cpu_usage', 'ts');
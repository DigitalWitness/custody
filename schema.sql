create table if not exists identities (
  id integer not null primary key,
  name text not null,
  created_at timestamp not null,
  public_key blob not null -- an x509 cert as ascii
);

create table if not exists ledger (
  id integer not null primary key,
  created_at timestamp not null,
  identity integer not null,
  message text not null,
  hash blob not null,

  foreign key (identity) references identities(id)
);

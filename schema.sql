create table if not exists identities (
  id integer not null primary key,
  name text not null,
  created_at timestamp not null,
  public_key blob not null -- an x509 cert as ascii
);

-- so we can look users up by their name
CREATE INDEX username_idx
  ON identities (name);

-- so we can look users up by their public key
CREATE INDEX publickey_idx
  ON identities (public_key);

create table if not exists ledger (
  id integer not null primary key,
  created_at timestamp not null,
  identity integer not null,
  message text not null,
  parent blob not null, -- signature of previous message
  signature blob not null, -- ecdsa signature of the message and parent fields

  foreign key (identity) references identities(id)
);

-- so we can find all messages from a user
CREATE INDEX ledger_identity_idx
  ON ledger (identity);

-- so we can sort all messages by timestamp
CREATE INDEX ledger_createdat_idx
  ON ledger (created_at);

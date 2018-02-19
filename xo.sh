sqlite3 'custody.sqlite' < schema.sql
xo -k field -o models file://custody.sqlite

echo "Initializing database"
cat sql/initdb.sql | sqlite3 nij.db
mkdir -p uploads
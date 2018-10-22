#!/usr/bin/env sh

chmod +x /bin/gopkg

echo /bin/gopkg -db.host $DB_HOST -db.store $USE_DB_STORE
/bin/gopkg -db.host $DB_HOST -db.store $USE_DB_STORE

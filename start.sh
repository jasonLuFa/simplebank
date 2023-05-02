#!/bin/sh
# exit immediately if a command returns a non-zero status
set -e

echo "run db migration"
# shellcheck source=/dev/null
. /app/app.env # same as `source /app/app.env`
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
# run the command given by the command line parameters in such a way that the current process is replaced by it
# from dockerfile we run /app/start.sh /app/main
# so in this point it will run /app/main to replace this shell process
exec "$@"

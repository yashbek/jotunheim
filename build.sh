set -xe

for dir in ./db/models/*/;
  do find $dir -type f -delete;
done

sqlc generate
go build -o ./bin/exist main.go
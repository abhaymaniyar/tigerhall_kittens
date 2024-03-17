echo "Loading $1"
for i in $(cat $1); do
  export $i
done

go run ./cmd/main.go
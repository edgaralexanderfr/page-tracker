build:
	go build -o bin/page-tracker . \
	&& GOOS=windows GOARCH=amd64 go build -o bin/page-tracker-amd64.exe . \
	&& GOOS=windows GOARCH=386 go build -o bin/page-tracker-386.exe .

exec: build
	./bin/page-tracker

run:
	clear && go run .

clean:
	rm bin/page-tracker
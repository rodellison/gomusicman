.PHONY: buildForAWS buildForOSX copyJson clean deploy gomodgen

buildForAWS: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/main main.go

buildForOSX: gomodgen
	export GO111MODULE=on
	env GOOS=darwin go build -ldflags="-s -w" -o bin/main main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

copyJson:
	mkdir bin
	cp apl_template_export.json bin/apl_template_export.json

testOSX: clean copyJson buildForOSX
	go test -v -covermode count ./...

testAWS: clean copyJson buildForAWS
	go test -v -covermode count ./...

deploy: clean copyJson buildForAWS
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh


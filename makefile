# todo add check to verify sbot is installed

build:
	go build -o semver-order cmd/semver-order/main.go

clean:
	rm -rf semver-order

test:
	go test ./...

release.patch:
	sbot release version -m patch
	sbot push version
	sbot get version

release.minor:
	sbot release version -m minor
	sbot push version
	sbot get version

release.major:
	sbot release version -m major
	sbot push version
	sbot get version
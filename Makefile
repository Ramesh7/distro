build: clean format
	go build

build-osx: clean format
	go build
	tar cvzf distro-v$(VERSION)-osx.tgz distro

build-linux: clean format
	env GOOS=linux go build
	tar cvzf distro-v$(VERSION)-linux.tgz distro

dist:
	rm -rf *.tgz dist/
	mkdir dist
	@$(MAKE) build-osx
	@$(MAKE) build-linux
	@$(MAKE) clean
	mv *.tgz dist/

format:
	go fmt

clean:
	go clean

.PHONY: build format clean dist build-osx build-linux

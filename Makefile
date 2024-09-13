TARGET := bin/configamajig
OSX-INTEL-TARGET := bin-osx-intel/configamajig
OSX-SILICONE-TARGET := bin-osx-silicone/configamajig
LINUX-TARGET := bin-linux/configamajig
WINDOWS-TARGET := bin-windows/configamajig

build: $(TARGET)

build-all: $(OSX-INTEL-TARGET) $(OSX-SILICONE-TARGET) $(LINUX-TARGET) $(WINDOWS-TARGET)

$(TARGET): * cmd/* actions/*
	go build -tags dynamic -o $@

$(OSX-INTEL-TARGET):
	env GOOS=darwin GOARCH=amd64 go build -o $@ 

$(OSX-SILICONE-TARGET):
	env GOOS=darwin GOARCH=arm64 go build -o $@ 

$(LINUX-TARGET):
	env GOOS=linux GOARCH=amd64 go build -o $@

$(WINDOWS-TARGET):
	env GOOS=windows GOARCH=amd64 go build -o $@

clean:
	go clean
	rm -rf $(TARGET)
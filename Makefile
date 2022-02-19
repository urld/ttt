PACKAGE      = ttt
VERSION      = $(shell git log -n1 --pretty='%h')
BUILD_DIR    = build
RELEASE_DIR  = dist
RELEASE_FILE = $(PACKAGE)_$(VERSION)_$(shell go env GOOS)-$(shell go env GOARCH)

.PHONY: all clean clean_build clean_dist dist build install test


all: install test



dist: build shrink
	mkdir -p $(RELEASE_DIR)
	mkdir -p $(BUILD_DIR)/licenses
	cp LICENSE $(BUILD_DIR)/licenses/ttt.LICENSE
	tar -cvzf  $(RELEASE_DIR)/$(RELEASE_FILE).tar.gz $(BUILD_DIR) --transform='s/$(BUILD_DIR)/$(RELEASE_FILE)/g'

shrink: build
	strip $(BUILD_DIR)/ttt*
	upx $(BUILD_DIR)/ttt*

build: clean_build
	mkdir -p $(BUILD_DIR)
	cd $(BUILD_DIR) && \
	CGO_ENABLED=1 go build -tags "sqlite_foreign_keys" github.com/urld/ttt/cmd/ttt

test:
	go test -v github.com/urld/ttt/...


install:
	go install github.com/urld/ttt/cmd/ttt


clean: clean_build clean_dist


clean_build:
	rm -rf $(BUILD_DIR)


clean_dist:
	rm -rf $(RELEASE_DIR)


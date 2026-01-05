# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names for examples
EXAMPLE_MAIN_LIST_DEVICES=main--devices-list
EXAMPLE_MAIN_DEVICE_STATUS=main-device-status
EXAMPLE_MAIN_DEVICE_INFO=main-device-info
EXAMPLE_MAIN_DEVICE_ERRORS=main-device-errors
EXAMPLE_RDWS_REBOOT=rdws-reboot
EXAMPLE_MAIN_DEVICE_LOCAL_DWS=main-device-local-dws
EXAMPLE_RDWS_SNAPSHOT=rdws-snapshot
EXAMPLE_RDWS_REPROVISION=rdws-reprovision
EXAMPLE_RDWS_DWS_PASSWORD=rdws-dws-password
EXAMPLE_MAIN_DEVICE_CHANGE_GROUP=main-device-change-group
EXAMPLE_MAIN_DEVICE_DELETE=main-device-delete
EXAMPLE_MAIN_GROUP_INFO=main-group-info
EXAMPLE_MAIN_GROUP_UPDATE=main-group-update
EXAMPLE_MAIN_GROUP_DELETE=main-group-delete
EXAMPLE_MAIN_LOCAL_DWS=main-local-dws
EXAMPLE_MAIN_AUTH_INFO=main-auth-info
EXAMPLE_MAIN_TEST_TOKEN=main-token-test
EXAMPLE_MAIN_TEST_ENDPOINTS=main-endpoints-test
EXAMPLE_RDWS_INFO=rdws-info
EXAMPLE_RDWS_TIME=rdws-time
EXAMPLE_RDWS_HEALTH=rdws-health
EXAMPLE_BDEPLOY_GET_SETUP=bdeploy-get-setup
EXAMPLE_BDEPLOY_DELETE_SETUP=bdeploy-delete-setup
EXAMPLE_BDEPLOY_DELETE_DEVICE=bdeploy-delete-device
EXAMPLE_BDEPLOY_GET_DEVICE=bdeploy-get-device
EXAMPLE_BDEPLOY_LIST_DEVICES=bdeploy-list-devices
EXAMPLE_BDEPLOY_ADD_SETUP=bdeploy-add-setup
EXAMPLE_BDEPLOY_UPDATE_SETUP=bdeploy-update-setup
EXAMPLE_BDEPLOY_LIST_SETUPS=bdeploy-list-setups
EXAMPLE_BDEPLOY_ASSOCIATE=bdeploy-associate
EXAMPLE_BDEPLOY_FIND_DEVICE=bdeploy-find-device
EXAMPLE_RDWS_FILES_LIST=rdws-files-list
EXAMPLE_RDWS_FILES_UPLOAD=rdws-files-upload
EXAMPLE_RDWS_FILES_RENAME=rdws-files-rename
EXAMPLE_RDWS_FILES_DELETE=rdws-files-delete
EXAMPLE_RDWS_LOGS_GET=rdws-logs-get
EXAMPLE_RDWS_CRASHDUMP_GET=rdws-crashdump-get
EXAMPLE_RDWS_FIRMWARE_DOWNLOAD=rdws-firmware-download
EXAMPLE_RDWS_REFORMAT_STORAGE=rdws-reformat-storage
EXAMPLE_RDWS_SSH=rdws-ssh
EXAMPLE_RDWS_TELNET=rdws-telnet
EXAMPLE_RDWS_REGISTRY_GET=rdws-registry-get
EXAMPLE_RDWS_REGISTRY_SET=rdws-registry-set
EXAMPLE_MAIN_CONTENT_LIST=main-content-list
EXAMPLE_MAIN_CONTENT_DELETE=main-content-delete
EXAMPLE_MAIN_CONTENT_UPLOAD=main-content-upload
EXAMPLE_MAIN_CONTENT_DOWNLOAD=main-content-download
EXAMPLE_MAIN_DEVICE_DOWNLOADS=main-device-downloads
EXAMPLE_MAIN_DEVICE_OPERATIONS=main-device-operations
EXAMPLE_MAIN_LIST_SUBSCRIPTIONS=main-subscriptions-list
EXAMPLE_MAIN_SUBSCRIPTION_COUNT=main-subscription-count
EXAMPLE_MAIN_SUBSCRIPTION_OPERATIONS=main-subscription-operations
EXAMPLE_MAIN_PRESENTATION_COUNT=main-presentation-count
EXAMPLE_MAIN_PRESENTATION_INFO=main-presentation-info
EXAMPLE_MAIN_PRESENTATION_LIST=main-presentation-list
EXAMPLE_MAIN_PRESENTATION_DELETE=main-presentation-delete
EXAMPLE_MAIN_PRESENTATION_INFO_BY_NAME=main-presentation-info-by-name
EXAMPLE_MAIN_PRESENTATION_CREATE=main-presentation-create
EXAMPLE_MAIN_PRESENTATION_UPDATE=main-presentation-update
EXAMPLE_MAIN_PRESENTATION_DELETE_BY_FILTER=main-presentation-delete-by-filter
EXAMPLE_MAIN_GET_REGTOKEN=main-get-regtoken
EXAMPLE_MAIN_DEVICE_FIND=main-device-find

# Build directory
BUILDDIR=bin

all: test build-examples

build: build-examples

build-examples: build-main--devices-list build-main-device-status build-main-device-info build-main-device-errors build-rdws-reboot build-main-device-local-dws build-rdws-snapshot build-rdws-reprovision build-rdws-dws-password build-main-device-change-group build-main-device-delete build-main-device-downloads build-main-device-operations build-main-group-info build-main-group-update build-main-group-delete build-main-local-dws build-main-auth-info build-main-token-test build-main-endpoints-test build-rdws-info build-rdws-time build-rdws-health build-bdeploy-get-setup build-bdeploy-delete-setup build-bdeploy-delete-device build-bdeploy-get-device build-bdeploy-list-devices build-bdeploy-add-setup build-bdeploy-update-setup build-bdeploy-list-setups build-bdeploy-associate build-bdeploy-find-device build-rdws-files-list build-rdws-files-upload build-rdws-files-rename build-rdws-files-delete build-rdws-logs-get build-rdws-crashdump-get build-rdws-firmware-download build-rdws-reformat-storage build-rdws-ssh build-rdws-telnet build-rdws-registry-get build-rdws-registry-set build-main-content-list build-main-content-delete build-main-content-upload build-main-content-download build-main-subscriptions-list build-main-subscription-count build-main-subscription-operations build-main-presentation-count build-main-presentation-info build-main-presentation-list build-main-presentation-delete build-main-presentation-info-by-name build-main-presentation-create build-main-presentation-update build-main-presentation-delete-by-filter build-main-get-regtoken build-main-device-find

build-main--devices-list:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_LIST_DEVICES) -v ./examples/main--devices-list

build-main-device-status:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_STATUS) -v ./examples/main-device-status

build-main-device-info:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_INFO) -v ./examples/main-device-info

build-main-device-errors:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_ERRORS) -v ./examples/main-device-errors

build-rdws-reboot:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_REBOOT) -v ./examples/rdws-reboot

build-main-device-local-dws:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_LOCAL_DWS) -v ./examples/main-device-local-dws

build-rdws-snapshot:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_SNAPSHOT) -v ./examples/rdws-snapshot

build-rdws-reprovision:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_REPROVISION) -v ./examples/rdws-reprovision

build-rdws-dws-password:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_DWS_PASSWORD) -v ./examples/rdws-dws-password

build-main-device-change-group:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_CHANGE_GROUP) -v ./examples/main-device-change-group

build-main-device-delete:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_DELETE) -v ./examples/main-device-delete

build-main-device-downloads:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_DOWNLOADS) -v ./examples/main-device-downloads

build-main-device-operations:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_OPERATIONS) -v ./examples/main-device-operations

build-main-group-info:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_GROUP_INFO) -v ./examples/main-group-info

build-main-group-update:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_GROUP_UPDATE) -v ./examples/main-group-update

build-main-group-delete:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_GROUP_DELETE) -v ./examples/main-group-delete

build-main-local-dws:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_LOCAL_DWS) -v ./examples/main-local-dws

build-main-auth-info:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_AUTH_INFO) -v ./examples/main-auth-info

build-main-token-test:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_TEST_TOKEN) -v ./examples/main-token-test

build-main-endpoints-test:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_TEST_ENDPOINTS) -v ./examples/main-endpoints-test

build-rdws-info:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_INFO) -v ./examples/rdws-info

build-rdws-time:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_TIME) -v ./examples/rdws-time

build-rdws-health:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_HEALTH) -v ./examples/rdws-health

build-bdeploy-get-setup:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_GET_SETUP) -v ./examples/bdeploy-get-setup

build-bdeploy-delete-setup:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_DELETE_SETUP) -v ./examples/bdeploy-delete-setup

build-bdeploy-delete-device:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_DELETE_DEVICE) -v ./examples/bdeploy-delete-device

build-bdeploy-get-device:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_GET_DEVICE) -v ./examples/bdeploy-get-device

build-bdeploy-list-devices:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_LIST_DEVICES) -v ./examples/bdeploy-list-devices

build-bdeploy-add-setup:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_ADD_SETUP) -v ./examples/bdeploy-add-setup

build-bdeploy-update-setup:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_UPDATE_SETUP) -v ./examples/bdeploy-update-setup

build-bdeploy-list-setups:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_LIST_SETUPS) -v ./examples/bdeploy-list-setups

build-bdeploy-associate:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_ASSOCIATE) -v ./examples/bdeploy-associate

build-bdeploy-find-device:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_BDEPLOY_FIND_DEVICE) -v ./examples/bdeploy-find-device

build-rdws-files-list:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_FILES_LIST) -v ./examples/rdws-files-list

build-rdws-files-upload:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_FILES_UPLOAD) -v ./examples/rdws-files-upload

build-rdws-files-rename:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_FILES_RENAME) -v ./examples/rdws-files-rename

build-rdws-files-delete:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_FILES_DELETE) -v ./examples/rdws-files-delete

build-rdws-logs-get:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_LOGS_GET) -v ./examples/rdws-logs-get

build-rdws-crashdump-get:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_CRASHDUMP_GET) -v ./examples/rdws-crashdump-get

build-rdws-firmware-download:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_FIRMWARE_DOWNLOAD) -v ./examples/rdws-firmware-download

build-rdws-reformat-storage:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_REFORMAT_STORAGE) -v ./examples/rdws-reformat-storage

build-rdws-ssh:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_SSH) -v ./examples/rdws-ssh

build-rdws-telnet:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_TELNET) -v ./examples/rdws-telnet

build-rdws-registry-get:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_REGISTRY_GET) -v ./examples/rdws-registry-get

build-rdws-registry-set:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_REGISTRY_SET) -v ./examples/rdws-registry-set

build-main-content-list:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_CONTENT_LIST) -v ./examples/main-content-list

build-main-content-delete:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_CONTENT_DELETE) -v ./examples/main-content-delete

build-main-content-upload:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_CONTENT_UPLOAD) -v ./examples/main-content-upload

build-main-content-download:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_CONTENT_DOWNLOAD) -v ./examples/main-content-download

build-main-subscriptions-list:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_LIST_SUBSCRIPTIONS) -v ./examples/main-subscriptions-list

build-main-subscription-count:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_SUBSCRIPTION_COUNT) -v ./examples/main-subscription-count

build-main-subscription-operations:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_SUBSCRIPTION_OPERATIONS) -v ./examples/main-subscription-operations

build-main-presentation-count:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_PRESENTATION_COUNT) -v ./examples/main-presentation-count

build-main-presentation-info:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_PRESENTATION_INFO) -v ./examples/main-presentation-info

build-main-presentation-list:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_PRESENTATION_LIST) -v ./examples/main-presentation-list

build-main-presentation-delete:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_PRESENTATION_DELETE) -v ./examples/main-presentation-delete

build-main-presentation-info-by-name:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_PRESENTATION_INFO_BY_NAME) -v ./examples/main-presentation-info-by-name

build-main-presentation-create:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_PRESENTATION_CREATE) -v ./examples/main-presentation-create

build-main-presentation-update:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_PRESENTATION_UPDATE) -v ./examples/main-presentation-update

build-main-presentation-delete-by-filter:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_PRESENTATION_DELETE_BY_FILTER) -v ./examples/main-presentation-delete-by-filter

build-main-get-regtoken:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_GET_REGTOKEN) -v ./examples/main-get-regtoken

build-main-device-find:
	mkdir -p $(BUILDDIR)
	$(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_FIND) -v ./examples/main-device-find

build-all-examples:
	mkdir -p $(BUILDDIR)
	for example in examples/*/; do \
		name=$$(basename "$$example"); \
		$(GOBUILD) -o $(BUILDDIR)/$$name -v ./examples/$$name; \
	done

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BUILDDIR)

run-main--devices-list: build-main--devices-list
	./$(BUILDDIR)/$(EXAMPLE_MAIN_LIST_DEVICES)

run-main-device-status: build-main-device-status
	./$(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_STATUS)

run-main-device-info: build-main-device-info
	./$(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_INFO)

run-main-device-errors: build-main-device-errors
	./$(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_ERRORS)

run-rdws-reboot: build-rdws-reboot
	./$(BUILDDIR)/$(EXAMPLE_RDWS_REBOOT)

run-main-device-local-dws: build-main-device-local-dws
	./$(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_LOCAL_DWS)

run-rdws-snapshot: build-rdws-snapshot
	./$(BUILDDIR)/$(EXAMPLE_RDWS_SNAPSHOT)

run-rdws-reprovision: build-rdws-reprovision
	./$(BUILDDIR)/$(EXAMPLE_RDWS_REPROVISION)

run-rdws-dws-password: build-rdws-dws-password
	./$(BUILDDIR)/$(EXAMPLE_RDWS_DWS_PASSWORD)

run-main-device-downloads: build-main-device-downloads
	./$(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_DOWNLOADS)

run-main-device-operations: build-main-device-operations
	./$(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_OPERATIONS)

run-main-local-dws: build-main-local-dws
	./$(BUILDDIR)/$(EXAMPLE_MAIN_LOCAL_DWS)

run-main-auth-info: build-main-auth-info
	./$(BUILDDIR)/$(EXAMPLE_MAIN_AUTH_INFO)

run-main-token-test: build-main-token-test
	./$(BUILDDIR)/$(EXAMPLE_MAIN_TEST_TOKEN)

run-rdws-info: build-rdws-info
	./$(BUILDDIR)/$(EXAMPLE_RDWS_INFO)

run-rdws-time: build-rdws-time
	./$(BUILDDIR)/$(EXAMPLE_RDWS_TIME)

run-rdws-health: build-rdws-health
	./$(BUILDDIR)/$(EXAMPLE_RDWS_HEALTH)

run-bdeploy-get-setup: build-bdeploy-get-setup
	./$(BUILDDIR)/$(EXAMPLE_BDEPLOY_GET_SETUP)

run-bdeploy-get-device: build-bdeploy-get-device
	./$(BUILDDIR)/$(EXAMPLE_BDEPLOY_GET_DEVICE)

run-bdeploy-list-devices: build-bdeploy-list-devices
	./$(BUILDDIR)/$(EXAMPLE_BDEPLOY_LIST_DEVICES)

run-rdws-files-list: build-rdws-files-list
	./$(BUILDDIR)/$(EXAMPLE_RDWS_FILES_LIST)

run-rdws-files-upload: build-rdws-files-upload
	./$(BUILDDIR)/$(EXAMPLE_RDWS_FILES_UPLOAD)

run-rdws-files-rename: build-rdws-files-rename
	./$(BUILDDIR)/$(EXAMPLE_RDWS_FILES_RENAME)

run-rdws-files-delete: build-rdws-files-delete
	./$(BUILDDIR)/$(EXAMPLE_RDWS_FILES_DELETE)

run-rdws-logs-get: build-rdws-logs-get
	./$(BUILDDIR)/$(EXAMPLE_RDWS_LOGS_GET)

run-rdws-crashdump-get: build-rdws-crashdump-get
	./$(BUILDDIR)/$(EXAMPLE_RDWS_CRASHDUMP_GET)

run-rdws-firmware-download: build-rdws-firmware-download
	./$(BUILDDIR)/$(EXAMPLE_RDWS_FIRMWARE_DOWNLOAD)

run-rdws-reformat-storage: build-rdws-reformat-storage
	./$(BUILDDIR)/$(EXAMPLE_RDWS_REFORMAT_STORAGE)

run-rdws-ssh: build-rdws-ssh
	./$(BUILDDIR)/$(EXAMPLE_RDWS_SSH)

run-rdws-telnet: build-rdws-telnet
	./$(BUILDDIR)/$(EXAMPLE_RDWS_TELNET)

run-rdws-registry-get: build-rdws-registry-get
	./$(BUILDDIR)/$(EXAMPLE_RDWS_REGISTRY_GET)

run-rdws-registry-set: build-rdws-registry-set
	./$(BUILDDIR)/$(EXAMPLE_RDWS_REGISTRY_SET)

run-main-content-list: build-main-content-list
	./$(BUILDDIR)/$(EXAMPLE_MAIN_CONTENT_LIST)

run-main-content-delete: build-main-content-delete
	./$(BUILDDIR)/$(EXAMPLE_MAIN_CONTENT_DELETE)

run-main-content-upload: build-main-content-upload
	./$(BUILDDIR)/$(EXAMPLE_MAIN_CONTENT_UPLOAD)

run-main-subscriptions-list: build-main-subscriptions-list
	./$(BUILDDIR)/$(EXAMPLE_MAIN_LIST_SUBSCRIPTIONS)

run-main-subscription-count: build-main-subscription-count
	./$(BUILDDIR)/$(EXAMPLE_MAIN_SUBSCRIPTION_COUNT)

run-main-subscription-operations: build-main-subscription-operations
	./$(BUILDDIR)/$(EXAMPLE_MAIN_SUBSCRIPTION_OPERATIONS)

run-main-get-regtoken: build-main-get-regtoken
	./$(BUILDDIR)/$(EXAMPLE_MAIN_GET_REGTOKEN)

run-main-device-find: build-main-device-find
	./$(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_FIND)

run: run-main--devices-list

run-examples: build-examples
	@echo "Available examples:"
	@ls -la $(BUILDDIR)/

deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Cross compilation
build-linux: build-linux-main--devices-list build-linux-main-device-status build-linux-main-device-errors build-linux-rdws-reboot build-linux-main-device-local-dws build-linux-rdws-snapshot build-linux-rdws-reprovision build-linux-rdws-dws-password build-linux-main-local-dws

build-linux-main--devices-list:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_LIST_DEVICES)_linux -v ./examples/main--devices-list

build-linux-main-device-status:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_STATUS)_linux -v ./examples/main-device-status

build-linux-main-device-errors:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_ERRORS)_linux -v ./examples/main-device-errors

build-linux-rdws-reboot:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_REBOOT)_linux -v ./examples/rdws-reboot

build-linux-main-device-local-dws:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_DEVICE_LOCAL_DWS)_linux -v ./examples/main-device-local-dws

build-linux-rdws-snapshot:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_SNAPSHOT)_linux -v ./examples/rdws-snapshot

build-linux-rdws-reprovision:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_REPROVISION)_linux -v ./examples/rdws-reprovision

build-linux-rdws-dws-password:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_RDWS_DWS_PASSWORD)_linux -v ./examples/rdws-dws-password

build-linux-main-local-dws:
	mkdir -p $(BUILDDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(EXAMPLE_MAIN_LOCAL_DWS)_linux -v ./examples/main-local-dws

# Development helpers
dev-deps:
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	golangci-lint run

fmt:
	goimports -w .
	$(GOCMD) fmt ./...

# Install for local use
install: install-main--devices-list install-main-device-status install-main-device-errors install-rdws-reboot install-main-device-local-dws install-rdws-snapshot install-rdws-reprovision install-rdws-dws-password install-main-local-dws install-main-auth-info

install-main--devices-list:
	$(GOCMD) install ./examples/main--devices-list

install-main-device-status:
	$(GOCMD) install ./examples/main-device-status

install-main-device-errors:
	$(GOCMD) install ./examples/main-device-errors

install-rdws-reboot:
	$(GOCMD) install ./examples/rdws-reboot

install-main-device-local-dws:
	$(GOCMD) install ./examples/main-device-local-dws

install-rdws-snapshot:
	$(GOCMD) install ./examples/rdws-snapshot

install-rdws-reprovision:
	$(GOCMD) install ./examples/rdws-reprovision

install-rdws-dws-password:
	$(GOCMD) install ./examples/rdws-dws-password

install-main-local-dws:
	$(GOCMD) install ./examples/main-local-dws

install-main-auth-info:
	$(GOCMD) install ./examples/main-auth-info

# Help target
help:
	@echo "Available targets:"
	@echo "  make all                    - Run tests and build all examples"
	@echo "  make build                  - Build all examples (alias for build-examples)"
	@echo "  make build-examples         - Build all example programs"
	@echo "  make build-main--devices-list - Build main--devices-list example"
	@echo "  make build-main-device-status - Build main-device-status example"
	@echo "  make build-main-device-errors - Build main-device-errors example"
	@echo "  make build-rdws-reboot      - Build rdws-reboot example"
	@echo "  make build-main-device-local-dws - Build main-device-local-dws example"
	@echo "  make build-rdws-snapshot    - Build rdws-snapshot example"
	@echo "  make build-rdws-reprovision - Build rdws-reprovision example"
	@echo "  make build-rdws-dws-password - Build rdws-dws-password example"
	@echo "  make build-main-local-dws   - Build main-local-dws example"
	@echo "  make build-main-auth-info   - Build main-auth-info example"
	@echo "  make build-main-token-test  - Build main-token-test example"
	@echo "  make build-bdeploy-get-setup - Build bdeploy-get-setup example"
	@echo "  make test                   - Run all tests"
	@echo "  make clean                  - Clean build artifacts"
	@echo "  make run                    - Run main--devices-list example"
	@echo "  make run-main--devices-list - Run main--devices-list example"
	@echo "  make run-main-device-status - Run main-device-status example"
	@echo "  make run-main-device-errors - Run main-device-errors example"
	@echo "  make run-rdws-reboot        - Run rdws-reboot example"
	@echo "  make run-main-device-local-dws - Run main-device-local-dws example"
	@echo "  make run-rdws-snapshot      - Run rdws-snapshot example"
	@echo "  make run-rdws-reprovision   - Run rdws-reprovision example"
	@echo "  make run-rdws-dws-password  - Run rdws-dws-password example"
	@echo "  make run-main-local-dws     - Run main-local-dws example"
	@echo "  make run-main-auth-info     - Run main-auth-info example"
	@echo "  make run-main-token-test    - Run main-token-test example"
	@echo "  make run-bdeploy-get-setup - Run bdeploy-get-setup example"
	@echo "  make run-bdeploy-get-device - Run bdeploy-get-device example"
	@echo "  make run-bdeploy-list-devices - Run bdeploy-list-devices example"
	@echo "  make deps                   - Download and tidy dependencies"
	@echo "  make build-linux            - Cross-compile all examples for Linux"
	@echo "  make install                - Install all examples to GOPATH/bin"
	@echo "  make fmt                    - Format code"
	@echo "  make lint                   - Run linter"
	@echo "  make help                   - Show this help message"

.PHONY: all build test clean run deps build-linux dev-deps lint fmt install \
	build-examples build-all-examples run-examples \
	build-main--devices-list build-main-device-status build-main-device-info build-main-device-errors build-rdws-reboot build-main-device-local-dws build-rdws-snapshot build-rdws-reprovision build-rdws-dws-password build-main-device-change-group build-main-device-delete build-main-device-downloads build-main-device-operations build-main-group-info build-main-group-update build-main-group-delete build-main-local-dws build-main-auth-info build-main-token-test build-main-endpoints-test build-rdws-info build-rdws-time build-rdws-health build-bdeploy-get-setup build-bdeploy-delete-setup build-bdeploy-delete-device build-bdeploy-get-device build-bdeploy-list-devices build-bdeploy-add-setup build-bdeploy-update-setup build-bdeploy-list-setups build-bdeploy-associate build-bdeploy-find-device build-rdws-files-list build-rdws-files-upload build-rdws-files-rename build-rdws-files-delete build-rdws-logs-get build-rdws-crashdump-get build-rdws-firmware-download build-rdws-reformat-storage build-rdws-ssh build-rdws-telnet build-rdws-registry-get build-rdws-registry-set build-main-content-list build-main-content-delete build-main-content-upload build-main-subscriptions-list build-main-subscription-count build-main-subscription-operations build-main-get-regtoken build-main-device-find \
	run-main--devices-list run-main-device-status run-main-device-info run-main-device-errors run-rdws-reboot run-main-device-local-dws run-rdws-snapshot run-rdws-reprovision run-rdws-dws-password run-main-device-downloads run-main-device-operations run-main-local-dws run-main-auth-info run-main-token-test run-rdws-info run-rdws-time run-rdws-health run-bdeploy-get-setup run-bdeploy-get-device run-bdeploy-list-devices run-rdws-files-list run-rdws-files-upload run-rdws-files-rename run-rdws-files-delete run-rdws-logs-get run-rdws-crashdump-get run-rdws-firmware-download run-rdws-reformat-storage run-rdws-ssh run-rdws-telnet run-rdws-registry-get run-rdws-registry-set run-main-content-list run-main-content-delete run-main-content-upload run-main-subscriptions-list run-main-subscription-count run-main-subscription-operations run-main-get-regtoken run-main-device-find \
	build-linux-main--devices-list build-linux-main-device-status build-linux-main-device-errors build-linux-rdws-reboot build-linux-main-device-local-dws build-linux-rdws-snapshot build-linux-rdws-reprovision build-linux-rdws-dws-password build-linux-main-local-dws \
	install-main--devices-list install-main-device-status install-main-device-errors install-rdws-reboot install-main-device-local-dws install-rdws-snapshot install-rdws-reprovision install-rdws-dws-password install-main-local-dws install-main-auth-info \
	help
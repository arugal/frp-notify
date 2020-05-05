# Copyright 2020 arugal, zhangwei24@apache.org
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSION ?= latest
OUT_DIR = bin
BINARY = frp-notify

OS = $(shell uname)

GO = go
GO_PATH = $$($(GO) env GOPATH)
GO_BUILD = $(GO) build
GO_GET = $(GO) get
GO_CLEAN = $(GO) clean
GO_TEST = $(GO) test
GO_INSTALL = $(GO) install
GO_LINT = $(GO_PATH)/bin/golangci-lint
GO_BUILD_FLAGS = -v
GO_BUILD_LDFLAGS = -X main.version=$(VERSION)

PLATFORMS := windows linux darwin

os = $(word 1, $@)
ARCH = amd64

.PHONY: deps
deps:
	$(GO_GET) -v -t -d ./...

.PHONE: build
build:
	${GO_BUILD} $(GO_BUILD_FLAGS) -o ${OUT_DIR}/${BINARY} cmd/main.go

.PHONE: build-all
build-all: windows linux darwin

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(os) GOARCH=$(ARCH) $(GO_BUILD) $(GO_BUILD_FLAGS) -ldflags "$(GO_BUILD_LDFLAGS)" -o $(OUT_DIR)/$(BINARY)-$(os)-$(ARCH) cmd/main.go


.PHONY: lint
lint:
	$(GO_LINT) version || curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GO_PATH}/bin v1.21.0
	$(GO_LINT) run --config ./golangci.yml

.PHONY: fix
fix:
	$(GO_LINT) run -v --fix ./...

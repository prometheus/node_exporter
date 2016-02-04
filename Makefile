# Copyright 2015 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSION  := 0.12.0rc3
TARGET   := node_exporter

REVISION := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
BRANCH   := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown')

REPO_PATH := "github.com/prometheus/node_exporter"
LDFLAGS   := -X main.Version=$(VERSION)
LDFLAGS   += -X $(REPO_PATH)/collector.Version=$(VERSION)
LDFLAGS   += -X $(REPO_PATH)/collector.Revision=$(REVISION)
LDFLAGS   += -X $(REPO_PATH)/collector.Branch=$(BRANCH)

GOFLAGS   := -ldflags "$(LDFLAGS)"

include Makefile.COMMON

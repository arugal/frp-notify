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

FROM golang:1.18-alpine AS builder

ENV GOPROXY="https://goproxy.cn"

ADD . /frp-notify
WORKDIR /frp-notify
RUN go mod tidy && \
    go build -o bin/frp-notify cmd/main.go

FROM alpine:latest

COPY --from=builder /frp-notify/bin/frp-notify .
ENTRYPOINT ["/frp-notify"]

CMD [ "start" ]
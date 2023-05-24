FROM scr.tiab.labda.sva.de/sthings-golang/sthings-golang:1.2 AS builder
LABEL maintainer="Patrick Hermann patrick.hermann@sva.de"

ARG VERSION=""
ARG BUILD_DATE=""
ARG COMMIT=""
ARG GIT_PAT=""

WORKDIR /src/
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/yacht-worker \
    -ldflags="-X main.version=v${VERSION} -X main.date=${BUILD_DATE} -X main.commit=${COMMIT}"

FROM alpine:3.17.0
COPY --from=builder /bin/yacht-worker /bin/yacht-worker

ENTRYPOINT ["yacht-worker"]
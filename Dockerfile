# Stage1: build from source
FROM golang:1.11.3 AS build

WORKDIR /work

COPY . /work

RUN CGO_ENABLED=0 GOLDFLAGS="-w -s" go build ./cmd/server

# Stage2: setup runtime container
FROM scratch

COPY --from=build /work/server /

EXPOSE 8080

ENTRYPOINT ["/server"]

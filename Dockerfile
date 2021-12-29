FROM golang:1.17 AS build-go
ADD . /src
WORKDIR /src
RUN go build -o regolint .

FROM gcr.io/distroless/base
COPY --from=build-go /src/regolint /regolint
WORKDIR /
ENTRYPOINT ["/regolint"]

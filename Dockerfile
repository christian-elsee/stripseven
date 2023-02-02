FROM docker.io/golang:1.18-alpine as build
LABEL stage=build
WORKDIR /usr/src/app
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN  go mod download && go mod verify

COPY . ./
RUN go build -v -o ./build main.go

FROM scratch
COPY --from=build /usr/src/app/build /build

ENTRYPOINT [ "/build" ]
CMD [ "-h" ]

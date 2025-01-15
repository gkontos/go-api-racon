#build executable 
FROM golang:1.22@sha256:88cf634d8c61055ab7c3e3e0026fcea2a3b5c9842c312c26ea4bd8e0bee2f7b1 as build-env

# git is installed to allow dependency installation from git sources
RUN apt update && apt install git
WORKDIR /app
COPY  . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/goapi

# create a minimal run image
FROM scratch

COPY --from=build-env /app/goapi /app/goapi

# the service listens on port 8443.
EXPOSE 443

CMD ["/app/goapi"]

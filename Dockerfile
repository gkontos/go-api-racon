#build executable 
FROM golang:1.22@sha256:a0679accac8685cc5389bd2298e045e570100940e6bdcca666a8ca7b32a1276c as build-env

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

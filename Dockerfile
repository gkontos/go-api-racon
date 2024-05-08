#build executable 
FROM golang:1.22@sha256:de2681281323110fcb09f0f520a22e77ac6ae3595c591add0982a55f0b8fa67b as build-env

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

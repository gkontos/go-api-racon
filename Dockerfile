#build executable 
FROM golang:1.22@sha256:2bd56f00ff47baf33e64eae7996b65846c7cb5e0a46e0a882ef179fd89654afa as build-env

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

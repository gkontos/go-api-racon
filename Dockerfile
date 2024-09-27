#build executable 
FROM golang:1.22@sha256:f8774f813f66d8c18679ad697f0e5d96a06f175795281cd1109eddf59b4c2419 as build-env

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

#build executable 
FROM golang:1.22@sha256:5370d4968adad7e969494e744c6d28a93931b89f259accf4d08a94c30446d3a3 as build-env

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

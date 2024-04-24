FROM golang:1.22.2 as build
#RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app
COPY . .
COPY api-weatherapi-com.pem /etc/ssl/certs
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/* /etc/ssl/certs/
COPY --from=build /app/cmd/cloudrun .
COPY --from=build /app/app.env .
ENTRYPOINT [ "./cloudrun" ]
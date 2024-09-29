FROM alpine:latest
WORKDIR /app
COPY loggerServiceApp /app/loggerServiceApp
CMD ["./loggerServiceApp"]

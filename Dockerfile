# Gebruik een minimale Go-builder image als basis
FROM golang:1.16-alpine AS builder

# Stel het werkdirectory in
WORKDIR /app

# Kopieer de broncode naar het werkdirectory
COPY . .

# Bouw de Go-toepassing
RUN go build -o main

# Gebruik een minimale alpine image als basis voor de runtime
FROM alpine:latest

# Installeer noodzakelijke pakketten
RUN apk --no-cache add ca-certificates

# Stel het werkdirectory in
WORKDIR /app

# Kopieer de gecompileerde binary vanuit de builder image
COPY --from=builder /app/main /app/main

# Kopieer de HTML-pagina naar het werkdirectory in de image
COPY helloworld.html /app/helloworld.html

# Exposeer de poort waarop de webtoepassing luistert
EXPOSE 8080

# Voer de Go-toepassing uit
CMD ["./main"]

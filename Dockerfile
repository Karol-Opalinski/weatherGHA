# ETAP 1 - budowa aplikacji bazowej
# obraz bazowy posiadający kompilator go
FROM golang:1.26-alpine AS builder

# ustawienie katalogu roboczego
WORKDIR /app

# skopiowanie skryptu do obrazu bazowego
COPY . .

# kompilacja aplikacji Go do postaci statycznego pliku wykonywalnego server
# CGO_ENABLED=0 wyłącza wykorzystanie bibliotek C
# -s -w usuwa informacje debug
RUN CGO_ENABLED=0 go build -o server -ldflags="-s -w" .

# ----------------------------------------
# ETAP 2 - konfiguracja serwera
# minimalny obraz bazowy scratch
FROM scratch

# Informacje OCI
LABEL org.opencontainers.image.authors="Karol Opaliński"

# ustawienie katalogu roboczego
WORKDIR /app

# skopiowanie kodu binarnego aplikacji oraz potrzebnych plików projektu
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/css ./css
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/cities.json .

# skopiowanie certyfikatu TLS dla połączeń https API
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# informacja o wystawianym porcie
EXPOSE 8080

# tryb healthcheck uruchamiany przez Docker HEALTHCHECK
# jeżeli program zostanie uruchomiony z argumentem "--healthcheck"
# proces kończy się natychmiast kodem 0, co oznacza poprawny stan aplikacji
HEALTHCHECK --interval=10s --timeout=3s \
CMD ["/app/server", "--healthcheck"]

# uruchomienie pliku binarnego aplikacji przy starcie kontenera
ENTRYPOINT ["/app/server"]
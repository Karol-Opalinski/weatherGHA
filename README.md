# Weather App – CI/CD (GitHub Actions + Docker)

## Opis projektu

Projekt przedstawia aplikację kontenerową z automatycznym pipeline CI/CD opartym o GitHub Actions. Pipeline buduje obraz Docker, wykonuje skan bezpieczeństwa Trivy i publikuje obraz do GitHub Container Registry (GHCR).

Dodatkowo wykorzystano cache z DockerHub oraz wsparcie dla wielu architektur (linux/amd64 i linux/arm64).

---

## CI/CD – główne kroki

Pipeline wykonuje:

- checkout kodu
- konfigurację Buildx i QEMU (multi-arch)
- logowanie do GHCR i DockerHub
- budowę obrazu Docker
- skan bezpieczeństwa Trivy
- publikację obrazu do GHCR

---

## Sekrety i uwierzytelnianie

W pipeline wykorzystano GitHub Secrets do przechowywania danych logowania:

- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

Do logowania w GHCR używany jest automatycznie dostępny token:

- `GITHUB_TOKEN`

Token ten działa poprawnie po ustawieniu w repozytorium opcji:
**Workflow permissions - Read and write permissions**

Dzięki temu GitHub może automatycznie publikować obrazy bez ręcznego tworzenia PAT.

---

## Trivy – problem i rozwiązanie

Podczas implementacji skanowania bezpieczeństwa wystąpił problem z użyciem eksportu obrazu do pliku `.tar`.

Początkowo obraz był zapisywany jako archiwum OCI i przekazywany do Trivy, jednak powodowało to błąd:

- brak pliku `manifest.json`
- Trivy nie rozpoznawał obrazu jako poprawnego archiwum Docker

Rozwiązanie polegało na zmianie podejścia:

- użycie `load: true` w Buildx

---

## Tagowanie obrazów

Obrazy są tagowane w dwóch wersjach:

- `latest` – najnowsza wersja
- `sha-<commit>` – wersja powiązana z konkretnym commit

---

## Cache

W pipeline zastosowano cache oparty o DockerHub

- `type=registry` – cache przechowywany w zewnętrznym rejestrze
- `mode=max` – maksymalna ilość warstw cache dla przyspieszenia kolejnych buildów

Cache jest zapisywany w dedykowanym publicznym repozytorium DockerHub i wykorzystywany w kolejnych uruchomieniach workflow w celu skrócenia czasu budowania obrazu.

---

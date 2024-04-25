# cloud-run

### A aplicação está disponível na seguinte url: https://cloud-run-goexpert-yt254bdq3q-uc.a.run.app/weather
### Exemplo de uso:
        curl -XPOST -H 'Content-Type: application/json' -d '{"zipcode":"71218010"}' 'https://cloud-run-goexpert-yt254bdq3q-uc.a.run.app/weather'

## Para rodar localmente, utilize o docker-compose.yaml:
        docker compose up -d
        docker exec -it gcloud-app bash
        go test ./...
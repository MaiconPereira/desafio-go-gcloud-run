# Desafio Go — Cloud Run

API em Go que recebe um CEP brasileiro de 8 dígitos, consulta a cidade no ViaCEP e retorna a temperatura atual via WeatherAPI em Celsius, Fahrenheit e Kelvin.

## URL Cloud Run

```text
https://desafio-go-gcloud-run-33238867170.southamerica-east1.run.app
```

Exemplo de uso:

```text
https://desafio-go-gcloud-run-33238867170.southamerica-east1.run.app/weather/01001000
```

## Contrato da API

### Sucesso — HTTP 200

```http
GET /weather/{cep}
```

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

Também aceita query string:

```http
GET /weather?zipcode=01001000
```

Endpoint alternativo compatível:

```http
GET /temperatures/{cep}
```

Health check:

```http
GET /health-check
```

### CEP inválido — HTTP 422

Formato inválido (não tem 8 dígitos ou contém letras):

```text
invalid zipcode
```

### CEP não encontrado — HTTP 404

Formato correto, mas não existe na base ViaCEP:

```text
can not find zipcode
```

## Variáveis de ambiente

| Variável | Obrigatória | Descrição |
|---|---|---|
| `PORT` | não | Porta HTTP (padrão: `8080`) |
| `WEATHER_API_KEY` | sim | Chave da [WeatherAPI](https://www.weatherapi.com/) |
| `VIACEP_BASE_URL` | não | URL base do ViaCEP (padrão já configurado) |
| `WEATHER_BASE_URL` | não | URL base da WeatherAPI (padrão já configurado) |

## Rodando localmente

```bash
cp .env.example .env
# Edite .env e preencha WEATHER_API_KEY
export $(grep -v '^#' .env | xargs)
go run ./cmd/server
```

Teste:

```bash
curl http://localhost:8080/weather/01001000
```

## Rodando com Docker

```bash
docker build -t desafio-go-gcloud-run .
docker run --rm -p 8080:8080 -e WEATHER_API_KEY=sua_chave desafio-go-gcloud-run
```

Teste:

```bash
curl http://localhost:8080/weather/01001000
```

## Rodando com Docker Compose

```bash
cp .env.example .env
# Edite .env e preencha WEATHER_API_KEY
docker compose up --build
```

## Testes

```bash
go test ./...
```

## Deploy no Google Cloud Run

### Pré-requisitos

- [Google Cloud CLI](https://cloud.google.com/sdk/docs/install) instalado e autenticado
- Projeto GCP criado

### Configuração inicial

```bash
gcloud config set project SEU_PROJECT_ID

gcloud services enable \
  run.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com
```

Crie o repositório Docker no Artifact Registry:

```bash
gcloud artifacts repositories create desafio-go \
  --repository-format=docker \
  --location=southamerica-east1 \
  --description="Imagens Docker do desafio Cloud Run"
```

Armazene a chave da WeatherAPI no Secret Manager:

```bash
echo -n "sua_chave_weatherapi" | gcloud secrets create weather-api-key --data-file=-
```

Permita que o Cloud Run leia o secret:

```bash
PROJECT_NUMBER=$(gcloud projects describe SEU_PROJECT_ID --format="value(projectNumber)")

gcloud secrets add-iam-policy-binding weather-api-key \
  --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

Permita que o Cloud Build faça deploy:

```bash
gcloud projects add-iam-policy-binding SEU_PROJECT_ID \
  --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding SEU_PROJECT_ID \
  --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding SEU_PROJECT_ID \
  --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"
```

### Deploy manual via Cloud Build

```bash
gcloud builds submit --config cloudbuild.yaml \
  --substitutions=_REGION=southamerica-east1,_REPOSITORY=desafio-go,_SERVICE_NAME=desafio-go-gcloud-run,_WEATHER_API_SECRET=weather-api-key
```

Ao final, o Cloud Build exibe a URL da aplicação.

### Trigger automático (CI/CD)

Configure um trigger no Cloud Build:

```text
Event: Push to a branch
Branch: ^main$
Configuration: Cloud Build configuration file
Cloud Build configuration file location: cloudbuild.yaml
```

Substitution variables:

```text
_REGION=southamerica-east1
_REPOSITORY=desafio-go
_SERVICE_NAME=desafio-go-gcloud-run
_WEATHER_API_SECRET=weather-api-key
```

Com essa configuração, cada push na branch `main` executa testes, build da imagem, push para o Artifact Registry e deploy no Cloud Run automaticamente.

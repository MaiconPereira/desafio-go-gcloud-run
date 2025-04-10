# Desafio 8 - Cloud Run

API em Go que recebe um CEP brasileiro de 8 dígitos, consulta a cidade no ViaCEP e retorna a temperatura atual usando WeatherAPI em Celsius, Fahrenheit e Kelvin.

O projeto foi preparado para rodar localmente, em Docker e no Google Cloud Run. A imagem final usa `scratch`, certificados HTTPS e usuário não-root.

## URL Cloud Run

```text
https://desafio-go-gcloud-run-esdras-santos-375082332631.southamerica-east1.run.app/weather/01001000
```

URL base:

```text
https://desafio-go-gcloud-run-esdras-santos-375082332631.southamerica-east1.run.app
```

## Contrato

### Sucesso

```http
GET /weather/01001000
```

Endpoint compatível:

```http
GET /temperatures/01001000
```

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

Também é possível consultar por query string:

```http
GET /weather?zipcode=01001000
```

Health check:

```http
GET /health-check
```

### Erros

CEP inválido:

```text
HTTP 422
invalid zipcode
```

CEP não encontrado:

```text
HTTP 404
can not find zipcode
```

## Variáveis de ambiente

```bash
PORT=8080
WEATHER_API_KEY=sua_chave_weatherapi
```

`WEATHER_API_KEY` é obrigatória para consultar a temperatura na WeatherAPI.

## Rodando localmente

```bash
cp .env.example .env
export $(grep -v '^#' .env | xargs)
go run ./cmd/server
```

Teste a API:

```bash
curl http://localhost:8080/weather/01001000
```

## Rodando com Docker

```bash
docker build -t desafio-8-cloud-run .
docker run --rm -p 8080:8080 -e WEATHER_API_KEY=sua_chave_weatherapi desafio-8-cloud-run
```

Teste a API:

```bash
curl http://localhost:8080/weather/01001000
```

## Rodando com Docker Compose

```bash
cp .env.example .env
export $(grep -v '^#' .env | xargs)
docker compose up --build
```

## Testes

```bash
go test ./...
```

## Deploy no Google Cloud Run

### Instalando o gcloud no Ubuntu

Se o comando `gcloud` não existir, instale o Google Cloud CLI:

```bash
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates gnupg curl

curl https://packages.cloud.google.com/apt/doc/apt-key.gpg \
  | sudo gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg

echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" \
  | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list

sudo apt-get update
sudo apt-get install -y google-cloud-cli
```

Faça login:

```bash
gcloud init
gcloud auth login
```

### Deploy com Cloud Build

Ative as APIs necessárias:

```bash
gcloud config set project SEU_PROJECT_ID
gcloud services enable run.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com secretmanager.googleapis.com
```

Crie o repositório Docker no Artifact Registry:

```bash
gcloud artifacts repositories create desafio-8 \
  --repository-format=docker \
  --location=us-central1 \
  --description="Imagens Docker do Desafio 8"
```

Crie o secret com a chave da WeatherAPI:

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

Permita que o Cloud Build faça deploy no Cloud Run:

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

Execute o pipeline:

```bash
gcloud builds submit --config cloudbuild.yaml \
  --substitutions=_REGION=us-central1,_REPOSITORY=desafio-8,_SERVICE_NAME=desafio-8-cloud-run,_WEATHER_API_SECRET=weather-api-key
```

Ao final, o Cloud Build mostra a URL publicada. Teste:

```bash
curl https://desafio-go-gcloud-run-esdras-santos-375082332631.southamerica-east1.run.app/weather/01001000
```

Resposta esperada:

```json
{
  "temp_C": 21,
  "temp_F": 69.8,
  "temp_K": 294.15
}
```

### Trigger pelo repositório

Se você criou um trigger no Cloud Build usando a opção **Dockerfile**, ele vai apenas construir a imagem. Para este projeto, use a opção **Cloud Build configuration file** apontando para:

```text
cloudbuild.yaml
```

Configuração recomendada do trigger:

```text
Event: Push to a branch
Branch: ^main$
Configuration: Cloud Build configuration file
Location: Repository
Cloud Build configuration file location: cloudbuild.yaml
```

Substitution variables do trigger:

```text
_REGION=us-central1
_REPOSITORY=desafio-8
_SERVICE_NAME=desafio-8-cloud-run
_WEATHER_API_SECRET=weather-api-key
```

Com essa configuração, cada push na branch `main` executa testes, build da imagem Docker, push para o Artifact Registry e deploy no Cloud Run.

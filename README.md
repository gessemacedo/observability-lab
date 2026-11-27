# observability-lab

Laboratório para praticar a **fase 1 do plano de carreira: observabilidade**.

A ideia é simples: subir uma API, expor métricas no formato Prometheus, ver essas métricas no Prometheus e acompanhar gráficos no Grafana. Depois, o mesmo serviço pode ir para um Kubernetes leve (k3s) no seu notebook.

Não precisa pagar nuvem. Tudo roda localmente com Docker Compose. O repositório é público para o GitHub Actions ficar no plano gratuito.

## O que tem aqui

- **API em Go** com três rotas:
  - `GET /` — hello em JSON
  - `GET /health` — liveness (`{"status":"ok"}`)
  - `GET /metrics` — métricas Prometheus (contador de requisições e histograma de latência por rota)
- **Dockerfile** em duas etapas: compilação em Go e imagem final mínima (distroless)
- **Docker Compose** com a API, Prometheus e Grafana já conectados
- **Dashboard do Grafana** provisionado: requisições/s e latência p50 / p90 / p99
- **Manifestos k3s/Kubernetes** em `k8s/`, com anotações de scrape do Prometheus
- **CI** no GitHub Actions: `go vet`, `go test`, `go build` e build da imagem Docker

## Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) e Docker Compose
- Opcional, para desenvolver a API fora do container: [Go 1.24+](https://go.dev/dl/)
- Opcional, para o exercício de Kubernetes: [k3s](https://k3s.io/) ou [k3d](https://k3d.io/) e `kubectl`

## Rodar localmente

Na raiz do repositório:

```bash
docker compose up --build
```

Isso sobe:

| Serviço    | URL                    | Para quê                         |
|------------|------------------------|----------------------------------|
| API        | http://localhost:8080  | Aplicação                        |
| Prometheus | http://localhost:9090  | Consultar e coletar métricas     |
| Grafana    | http://localhost:3000  | Gráficos já provisionados        |

Confira se a API está no ar:

```bash
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

Para desligar: `Ctrl+C` e, se quiser remover os containers, `docker compose down`.

### Grafana

Abra http://localhost:3000. O dashboard **Observability Lab** já entra sozinho, com o Prometheus como fonte de dados.

O acesso anônimo está ligado no perfil Viewer, então dá para ver os gráficos sem login. Se quiser entrar como admin, o usuário e a senha locais são `admin` / `admin` (só para este laboratório, não use isso em lugar nenhum de verdade).

### Gerar tráfego para os gráficos

Nos primeiros segundos os painéis ficam vazios: o Prometheus precisa coletar alguns intervalos. Gere algumas requisições e espere uns 15 segundos:

```bash
for i in $(seq 1 50); do
  curl -s http://localhost:8080/ >/dev/null
  curl -s http://localhost:8080/health >/dev/null
done
```

### Ver as métricas no Prometheus

1. Abra http://localhost:9090.
2. Em **Status → Targets**, o alvo `observability-lab` (`api:8080`) deve aparecer como `UP`.
3. Na caixa de consulta, experimente:

```promql
sum(rate(http_requests_total[1m])) by (path)
```

```promql
histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
```

A configuração de coleta está em `prometheus.yml` e já aponta para a API no Compose (`api:8080`).

## Desenvolver a API sem Docker

```bash
go test ./...
go vet ./...
go run .
```

A API sobe em `:8080`. Dá para mudar a porta com `PORT=9091 go run .`.

## Subir no k3s

O k3s, sozinho, não traz Prometheus nem Grafana. Este exercício sobe só a API no cluster. Para ver métricas, continue usando o Compose ou instale um stack de monitoramento (por exemplo, kube-prometheus-stack) que respeite as anotações `prometheus.io/scrape`.

```bash
docker build -t observability-lab:local .
```

Importe a imagem para o containerd do k3s (no Linux, com k3s instalado na máquina):

```bash
docker save observability-lab:local | sudo k3s ctr images import -
kubectl apply -f k8s/
kubectl rollout status deployment/observability-lab
kubectl port-forward svc/observability-lab 8080:8080
```

Depois, `curl http://localhost:8080/health` deve responder `{"status":"ok"}`.

Se você usa [k3d](https://k3d.io/), importe a imagem com:

```bash
k3d image import observability-lab:local
kubectl apply -f k8s/
```

## O que cada parte ensina

Este laboratório cobre a base da observabilidade, o primeiro degrau do plano:

- **Instrumentação.** A aplicação expõe o que acontece com ela. Não é “mais máquina”: é conseguir ver requisições, erros e latência por rota.
- **Prometheus.** Um coletor puxa métricas em intervalos curtos e deixa você perguntar coisas como “quantas requisições por segundo?” ou “qual a latência p99?”.
- **Grafana.** Métrica crua vira gráfico. Os percentis p50, p90 e p99 mostram o comportamento típico e a cauda lenta, que a média esconde.
- **Sinais de saúde.** `/health` é o que o Kubernetes usa para saber se o processo está vivo. É o começo de liveness e readiness.
- **CI.** Cada push ou pull request roda `go vet`, testes e o build da imagem. Você vê o resultado no GitHub Actions, sem publicar nem fazer deploy.
- **k3s.** Empacotar a API e subir num Kubernetes leve no notebook, com as anotações que um Prometheus no cluster usa para descobrir o alvo.

Próximos passos naturais (ainda neste laboratório, se quiser evoluir): traces (OpenTelemetry), logs estruturados, alertas no Prometheus e um dashboard por serviço em um cluster com Prometheus Operator.

## Licença e uso

Material de estudo. Não contém dados corporativos, segredos ou credenciais reais. A senha `admin` do Grafana vale só para o ambiente local deste repositório.

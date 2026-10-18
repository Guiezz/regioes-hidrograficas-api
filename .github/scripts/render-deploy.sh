#!/usr/bin/env bash
# Dispara um deploy do commit COMMIT_SHA no Render, acompanha até `live` e
# roda o teste de fumaça em SMOKE_URL.
#
# Obrigatórias: RENDER_API_KEY, RENDER_SERVICE_ID, COMMIT_SHA, SMOKE_URL.
# Opcionais: RENDER_API_BASE (https://api.render.com/v1), POLL_INTERVAL (15),
#            DEPLOY_TIMEOUT (1200), SMOKE_TIMEOUT (300), GITHUB_STEP_SUMMARY.
# Exit 0 = deploy live + fumaça ok; exit 1 = qualquer falha.
set -euo pipefail

RENDER_API_BASE="${RENDER_API_BASE:-https://api.render.com/v1}"
POLL_INTERVAL="${POLL_INTERVAL:-15}"
DEPLOY_TIMEOUT="${DEPLOY_TIMEOUT:-1200}"
SMOKE_TIMEOUT="${SMOKE_TIMEOUT:-300}"

err() { echo "❌ $*" >&2; }
info() { echo "$*" >&2; }

for var in RENDER_API_KEY RENDER_SERVICE_ID COMMIT_SHA SMOKE_URL; do
  if [[ -z "${!var:-}" ]]; then
    err "variável obrigatória ausente: $var"
    exit 1
  fi
done

TMP="$(mktemp -d)"
DEPLOY_ID="-"
DEPLOY_STATUS="não iniciado"
SMOKE_RESULT="não executado"

finish() {
  local code=$?
  if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
    {
      echo "### Deploy no Render"
      echo ""
      echo "| | |"
      echo "|---|---|"
      echo "| Commit | \`$COMMIT_SHA\` |"
      echo "| Deploy | \`$DEPLOY_ID\` |"
      echo "| Status final | $DEPLOY_STATUS |"
      echo "| Teste de fumaça | $SMOKE_RESULT |"
    } >> "$GITHUB_STEP_SUMMARY"
  fi
  rm -rf "$TMP"
  exit "$code"
}
trap finish EXIT

# render_api <método> <caminho> [corpo]: grava o corpo da resposta em
# $TMP/body e imprime o código HTTP. A chave só vai no header, nunca na saída.
render_api() {
  local method="$1" path="$2" data="${3:-}"
  local args=(-sS -o "$TMP/body" -w '%{http_code}' --max-time 30 -X "$method"
    -H "Authorization: Bearer $RENDER_API_KEY"
    -H "Accept: application/json")
  if [[ -n "$data" ]]; then
    args+=(-H "Content-Type: application/json" --data "$data")
  fi
  curl "${args[@]}" "$RENDER_API_BASE$path" || true
}

# 1. Dispara o deploy do commit.
info "Disparando deploy do commit $COMMIT_SHA no Render..."
code="$(render_api POST "/services/$RENDER_SERVICE_ID/deploys" "{\"commitId\":\"$COMMIT_SHA\"}")"
if ! DEPLOY_ID="$(jq -er '.id' "$TMP/body" 2>/dev/null)"; then
  DEPLOY_ID="-"
  DEPLOY_STATUS="falha ao disparar (HTTP $code)"
  err "resposta inesperada do Render ao disparar o deploy (HTTP $code):"
  cat "$TMP/body" >&2 2>/dev/null || true
  echo >&2
  exit 1
fi
info "Deploy criado: $DEPLOY_ID"

# 2. Acompanha até live, falha ou timeout.
start=$SECONDS
last=""
while true; do
  code="$(render_api GET "/services/$RENDER_SERVICE_ID/deploys/$DEPLOY_ID")"
  if ! status="$(jq -er '.status' "$TMP/body" 2>/dev/null)"; then
    DEPLOY_STATUS="resposta inesperada (HTTP $code)"
    err "resposta inesperada do Render ao consultar o deploy $DEPLOY_ID (HTTP $code):"
    cat "$TMP/body" >&2 2>/dev/null || true
    echo >&2
    exit 1
  fi
  DEPLOY_STATUS="$status"
  if [[ "$status" != "$last" ]]; then
    info "Status do deploy: $status"
    last="$status"
  fi
  case "$status" in
    live) break ;;
    build_failed | update_failed | canceled | deactivated | pre_deploy_failed)
      err "deploy $DEPLOY_ID terminou com status $status"
      exit 1
      ;;
  esac
  if (( SECONDS - start >= DEPLOY_TIMEOUT )); then
    DEPLOY_STATUS="timeout ($status)"
    err "timeout: deploy $DEPLOY_ID não ficou live em ${DEPLOY_TIMEOUT}s (último status: $status)"
    exit 1
  fi
  sleep "$POLL_INTERVAL"
done

# 3. Teste de fumaça: 200 com lista não vazia, tolerando o Render acordando.
smoke_interval=10
(( POLL_INTERVAL < smoke_interval )) && smoke_interval=$POLL_INTERVAL
info "Teste de fumaça em $SMOKE_URL..."
start=$SECONDS
while true; do
  code="$(curl -sS -o "$TMP/smoke" -w '%{http_code}' --max-time 30 "$SMOKE_URL" 2>/dev/null || true)"
  if [[ "$code" == "200" ]] && jq -e 'type == "array" and length > 0' "$TMP/smoke" >/dev/null 2>&1; then
    SMOKE_RESULT="ok (HTTP 200, lista não vazia)"
    info "✅ Teste de fumaça ok."
    break
  fi
  if (( SECONDS - start >= SMOKE_TIMEOUT )); then
    SMOKE_RESULT="falhou (último HTTP $code)"
    err "teste de fumaça falhou após ${SMOKE_TIMEOUT}s: último HTTP $code, esperado 200 com lista não vazia"
    exit 1
  fi
  info "Fumaça ainda não ok (HTTP $code), tentando de novo em ${smoke_interval}s..."
  sleep "$smoke_interval"
done

info "✅ Deploy $DEPLOY_ID live e respondendo."

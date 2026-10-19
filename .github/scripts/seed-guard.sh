#!/usr/bin/env bash
# Travas do workflow de seed.
#   seed-guard.sh confirm <valor>     exit 0 só se <valor> for exatamente SEED
#   seed-guard.sh check-log <arquivo> exit 1 se o log do seeder tiver linhas com ❌
#   seed-guard.sh clean-url < url     limpa a URL do banco (pontas, aspas) e valida o esquema
set -euo pipefail

case "${1:-}" in
  confirm)
    if [[ "${2-}" != "SEED" ]]; then
      echo "Confirmação inválida: digite SEED para rodar o seed." >&2
      exit 1
    fi
    echo "Confirmação ok."
    ;;
  check-log)
    log="${2:?uso: seed-guard.sh check-log <arquivo>}"
    if grep -F '❌' "$log" >&2; then
      echo "O seeder registrou erros (linhas com ❌ acima)." >&2
      exit 1
    fi
    echo "Log do seeder sem erros."
    ;;
  clean-url)
    # Lê a URL do stdin e a escreve limpa no stdout (sem quebra de linha no
    # fim). Nunca escreve a URL no stderr.
    url="$(cat)"
    trim() {
      local s="$1"
      s="${s#"${s%%[![:space:]]*}"}"
      printf '%s' "${s%"${s##*[![:space:]]}"}"
    }
    url="$(trim "$url")"
    if [[ ${#url} -ge 2 && "${url:0:1}" == "${url: -1}" && ( "${url:0:1}" == '"' || "${url:0:1}" == "'" ) ]]; then
      url="$(trim "${url:1:${#url}-2}")"
    fi
    # Espaço ou caractere de controle no meio também invalida: a URL vai para
    # o GITHUB_ENV, onde uma quebra de linha criaria outra variável.
    if [[ "$url" != postgres://* && "$url" != postgresql://* ]] || [[ "$url" =~ [[:space:][:cntrl:]] ]]; then
      echo "NEON_DATABASE_URL vazio ou inválido (esperado postgres:// ou postgresql://)" >&2
      exit 1
    fi
    printf '%s' "$url"
    ;;
  *)
    echo "uso: seed-guard.sh confirm <valor> | check-log <arquivo> | clean-url < url" >&2
    exit 2
    ;;
esac

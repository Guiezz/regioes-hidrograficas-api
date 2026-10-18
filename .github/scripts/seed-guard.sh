#!/usr/bin/env bash
# Travas do workflow de seed.
#   seed-guard.sh confirm <valor>     exit 0 só se <valor> for exatamente SEED
#   seed-guard.sh check-log <arquivo> exit 1 se o log do seeder tiver linhas com ❌
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
  *)
    echo "uso: seed-guard.sh confirm <valor> | check-log <arquivo>" >&2
    exit 2
    ;;
esac

#!/usr/bin/env bash
# Testa seed-guard.sh (confirmação digitada e checagem do log do seeder).
# Uso: .github/scripts/seed-guard.test.sh
set -uo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GUARD="$DIR/seed-guard.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

FAILED=0

# check <nome> <exit esperado> <texto esperado na saída ou ""> -- <args do guard...>
check() {
  local name="$1" want_exit="$2" want_out="$3"
  shift 4
  bash "$GUARD" "$@" > "$TMP/out" 2>&1
  local got_exit=$?

  local why=""
  [[ "$got_exit" == "$want_exit" ]] || why="exit $got_exit (esperado $want_exit)"
  if [[ -n "$want_out" ]] && ! grep -qF -- "$want_out" "$TMP/out"; then
    why="${why:+$why; }saída sem '$want_out'"
  fi

  if [[ -z "$why" ]]; then
    echo "PASS $name"
  else
    echo "FAIL $name: $why"
    sed 's/^/    saída: /' "$TMP/out"
    FAILED=1
  fi
}

check "confirm SEED" 0 "" -- confirm SEED
check "confirm seed" 1 "Confirmação inválida: digite SEED para rodar o seed." -- confirm seed
check "confirm vazio" 1 "Confirmação inválida: digite SEED para rodar o seed." -- confirm ""
check "confirm ' SEED'" 1 "Confirmação inválida: digite SEED para rodar o seed." -- confirm " SEED"

printf '🌊 PROCESSANDO BACIA: Alto Jaguaribe\n🚀 SEED COMPLETO\n' > "$TMP/ok.log"
check "check-log sem erro" 0 "" -- check-log "$TMP/ok.log"

printf '🌊 PROCESSANDO BACIA: X\n❌ Erro ao criar bacia X\n🚀 SEED COMPLETO\n' > "$TMP/erro.log"
check "check-log com ❌" 1 "❌ Erro ao criar bacia X" -- check-log "$TMP/erro.log"

exit "$FAILED"

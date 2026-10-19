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

URL="postgresql://dono:senha-de-teste@host.exemplo/db?sslmode=require&channel_binding=require"
URL_ERRO="NEON_DATABASE_URL vazio ou inválido (esperado postgres:// ou postgresql://)"

# check_clean <nome> <entrada no stdin> <exit esperado> <stdout esperado (exato)>
# Em qualquer caso, o stderr nunca pode conter a entrada nem a URL; em falha,
# precisa conter URL_ERRO.
check_clean() {
  local name="$1" input="$2" want_exit="$3" want_out="$4"
  printf '%s' "$input" | bash "$GUARD" clean-url > "$TMP/out" 2> "$TMP/err"
  local got_exit=$?

  local why=""
  [[ "$got_exit" == "$want_exit" ]] || why="exit $got_exit (esperado $want_exit)"
  if [[ "$want_exit" == 0 ]]; then
    [[ "$(cat "$TMP/out"; echo x)" == "${want_out}x" ]] || why="${why:+$why; }stdout diferente do esperado"
  else
    [[ -s "$TMP/out" ]] && why="${why:+$why; }stdout deveria estar vazio"
    grep -qF -- "$URL_ERRO" "$TMP/err" || why="${why:+$why; }stderr sem a mensagem de erro"
  fi
  local trimmed="${input//[[:space:]]/}"
  if [[ -n "$trimmed" ]] && grep -qF -- "$trimmed" "$TMP/err"; then
    why="${why:+$why; }stderr ecoou o valor"
  fi
  grep -qF -- "senha-de-teste" "$TMP/err" && why="${why:+$why; }stderr vazou a URL"

  if [[ -z "$why" ]]; then
    echo "PASS $name"
  else
    echo "FAIL $name: $why"
    sed 's/^/    stderr: /' "$TMP/err"
    FAILED=1
  fi
}

check_clean "clean-url quebra de linha no fim" "$URL"$'\n' 0 "$URL"
check_clean "clean-url CRLF no fim" "$URL"$'\r\n' 0 "$URL"
check_clean "clean-url espaços nas pontas" "  $URL  " 0 "$URL"
check_clean "clean-url aspas duplas" "\"$URL\"" 0 "$URL"
check_clean "clean-url aspas simples" "'$URL'" 0 "$URL"
check_clean "clean-url já limpa" "$URL" 0 "$URL"
check_clean "clean-url postgres://" "postgres://u:p@h/db" 0 "postgres://u:p@h/db"
check_clean "clean-url vazio" "" 1 ""
check_clean "clean-url só espaços" $'  \n' 1 ""
check_clean "clean-url mysql://x" "mysql://x" 1 ""
check_clean "clean-url quebra de linha no meio" "postgresql://a"$'\n'"b@h/db" 1 ""

exit "$FAILED"

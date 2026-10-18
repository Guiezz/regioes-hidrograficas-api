#!/usr/bin/env bash
# Testa render-deploy.sh contra um mock local da API do Render.
# Uso: .github/scripts/render-deploy.test.sh
set -uo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT="$DIR/render-deploy.sh"
TMP="$(mktemp -d)"
MOCK_PID=""

cleanup() {
  [[ -n "$MOCK_PID" ]] && kill "$MOCK_PID" 2>/dev/null
  rm -rf "$TMP"
}
trap cleanup EXIT

TEST_KEY="rnd_chave-de-teste-nao-pode-vazar"
# Sem "abc" no SHA: o cenário service_id_invalido exige que "abc" não apareça no stderr.
TEST_SERVICE="srv-abc123"
TEST_SHA="0123456789def"

cat > "$TMP/mock.py" <<'PY'
import json, os, sys
from http.server import BaseHTTPRequestHandler, HTTPServer

cfg = json.loads(os.environ["MOCK_CONFIG"])
key = os.environ["MOCK_KEY"]
service = os.environ["MOCK_SERVICE"]
sha = os.environ["MOCK_SHA"]
statuses = list(cfg.get("statuses", ["live"]))
smoke = list(cfg.get("smoke", [[200, [{"id": 1}]]]))

def next_of(seq):
    return seq.pop(0) if len(seq) > 1 else seq[0]

class H(BaseHTTPRequestHandler):
    def log_message(self, *a):
        pass

    def send(self, code, body):
        data = json.dumps(body).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def authorized(self):
        return self.headers.get("Authorization") == "Bearer " + key

    def do_POST(self):
        if self.path != f"/services/{service}/deploys":
            return self.send(404, {"message": "not found"})
        if "post" in cfg:
            return self.send(cfg["post"][0], cfg["post"][1])
        if not self.authorized():
            return self.send(401, {"message": "unauthorized"})
        body = json.loads(self.rfile.read(int(self.headers.get("Content-Length", 0))) or b"{}")
        if body.get("commitId") != sha:
            return self.send(400, {"message": "commitId errado"})
        self.send(201, {"id": "dep-test", "status": "created"})

    def do_GET(self):
        if self.path == "/smoke":
            code, body = next_of(smoke)
            return self.send(code, body)
        if self.path == f"/services/{service}/deploys/dep-test":
            if not self.authorized():
                return self.send(401, {"message": "unauthorized"})
            st = next_of(statuses)
            # Um item [código, corpo] simula uma resposta transitória do Render.
            if isinstance(st, list):
                return self.send(st[0], st[1])
            return self.send(200, {"id": "dep-test", "status": st})
        self.send(404, {"message": "not found"})

srv = HTTPServer(("127.0.0.1", 0), H)
with open(sys.argv[1], "w") as f:
    f.write(str(srv.server_address[1]))
srv.serve_forever()
PY

start_mock() {
  rm -f "$TMP/port"
  MOCK_CONFIG="$1" MOCK_KEY="$TEST_KEY" MOCK_SERVICE="$TEST_SERVICE" MOCK_SHA="$TEST_SHA" \
    python3 "$TMP/mock.py" "$TMP/port" &
  MOCK_PID=$!
  for _ in $(seq 1 50); do
    [[ -s "$TMP/port" ]] && break
    sleep 0.1
  done
  PORT="$(cat "$TMP/port")"
}

stop_mock() {
  kill "$MOCK_PID" 2>/dev/null
  wait "$MOCK_PID" 2>/dev/null
  MOCK_PID=""
}

FAILED=0

# run_case <nome> <config-json> <exit esperado> <texto esperado no stderr ou ""> [VAR=valor ...]
# FORBID_ERR=<texto> antes da chamada exige que o texto NÃO apareça no stderr.
run_case() {
  local name="$1" config="$2" want_exit="$3" want_err="$4"
  shift 4
  start_mock "$config"
  : > "$TMP/summary"
  env RENDER_API_KEY="$TEST_KEY" RENDER_SERVICE_ID="$TEST_SERVICE" COMMIT_SHA="$TEST_SHA" \
    RENDER_API_BASE="http://127.0.0.1:$PORT" SMOKE_URL="http://127.0.0.1:$PORT/smoke" \
    POLL_INTERVAL=0 DEPLOY_TIMEOUT=10 SMOKE_TIMEOUT=3 GITHUB_STEP_SUMMARY="$TMP/summary" \
    "$@" bash "$SCRIPT" > "$TMP/out" 2> "$TMP/err"
  local got_exit=$?
  stop_mock

  local why=""
  [[ "$got_exit" == "$want_exit" ]] || why="exit $got_exit (esperado $want_exit)"
  if [[ -n "$want_err" ]] && ! grep -q -- "$want_err" "$TMP/err"; then
    why="${why:+$why; }stderr sem '$want_err'"
  fi
  if grep -q -- "$TEST_KEY" "$TMP/out" "$TMP/err" "$TMP/summary"; then
    why="${why:+$why; }a chave da API vazou na saída"
  fi
  if grep -q -- "$TEST_SERVICE" "$TMP/out" "$TMP/err" "$TMP/summary"; then
    why="${why:+$why; }o service ID vazou na saída"
  fi
  if [[ -n "${FORBID_ERR:-}" ]] && grep -q -- "$FORBID_ERR" "$TMP/err"; then
    why="${why:+$why; }stderr contém '$FORBID_ERR'"
  fi
  if [[ "$got_exit" == 0 ]] && ! grep -q "dep-test" "$TMP/summary"; then
    why="${why:+$why; }GITHUB_STEP_SUMMARY sem o id do deploy"
  fi

  if [[ -z "$why" ]]; then
    echo "PASS $name"
  else
    echo "FAIL $name: $why"
    sed 's/^/    stderr: /' "$TMP/err"
    FAILED=1
  fi
}

run_case ok \
  '{"statuses":["build_in_progress","update_in_progress","live"],"smoke":[[200,[{"id":1}]]]}' \
  0 ""

run_case build_failed \
  '{"statuses":["build_in_progress","build_failed"]}' \
  1 "build_failed"

run_case 401 \
  '{"post":[401,{"message":"unauthorized"}]}' \
  1 "401"

run_case fumaca_lenta \
  '{"statuses":["live"],"smoke":[[503,{"message":"waking"}],[503,{"message":"waking"}],[200,[{"id":1}]]]}' \
  0 ""

run_case fumaca_vazia \
  '{"statuses":["live"],"smoke":[[200,[]]]}' \
  1 ""

run_case timeout \
  '{"statuses":["build_in_progress"]}' \
  1 "timeout" DEPLOY_TIMEOUT=1

run_case service_id_url \
  '{"statuses":["live"]}' \
  0 "" RENDER_SERVICE_ID="https://dashboard.render.com/web/srv-abc123/settings"

run_case service_id_espacos \
  '{"statuses":["live"]}' \
  0 "" RENDER_SERVICE_ID=" srv-abc123
"

FORBID_ERR=abc run_case service_id_invalido \
  '{"statuses":["live"]}' \
  1 "RENDER_SERVICE_ID inválido" RENDER_SERVICE_ID="abc"

run_case poll_transitorio \
  '{"statuses":[[502,"bad gateway"],[502,"bad gateway"],"build_in_progress","live"]}' \
  0 ""

run_case poll_persistente \
  '{"statuses":[[502,"bad gateway"],[502,"bad gateway"],[502,"bad gateway"],[502,"bad gateway"],"live"]}' \
  1 "HTTP 502"

exit "$FAILED"

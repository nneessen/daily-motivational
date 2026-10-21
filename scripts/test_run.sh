#!/usr/bin/env bash
# Checks that the app builds and actually runs. Run it after every change.
#   scripts/test_run.sh          format, vet, tests, build, then a live dry run (prints the email, sends nothing)
#   scripts/test_run.sh --send   same checks, then sends one real email instead of the dry run
set -euo pipefail
cd "$(dirname "$0")/.."

if [ -f .env ]; then
  set -a
  . ./.env
  set +a
fi

unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
  echo "FAIL: run gofmt -w on: $unformatted"
  exit 1
fi
go vet ./...
go test ./...
go build -o /dev/null .
echo "OK: format, vet, tests, build"

if [ -z "${ANTHROPIC_API_KEY:-}" ]; then
  echo "SKIPPED live check: set ANTHROPIC_API_KEY in .env to run it."
  exit 0
fi

if [ "${1:-}" = "--send" ]; then
  go run .
  echo "OK: email sent to ${TO_ADDRESS:-?}"
  exit 0
fi

out=$(go run . -dry-run)
echo "$out"
count=$(grep -c '^OPTION ' <<<"$out" || true)
if [ "$count" -ne 4 ]; then
  echo "FAIL: expected 4 options, got $count"
  exit 1
fi
echo "OK: live dry run returned 4 options"

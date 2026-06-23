#!/usr/bin/env bash
set -euo pipefail

file=$(jq -r '.tool_input.file_path // empty')
[[ "$file" =~ \.go$ ]] || exit 0

dir=$(dirname "$file")
out=$(golangci-lint run "${dir}/..." 2>&1) || true
[ -z "$out" ] && exit 0

python3 -c "
import json, sys
data = sys.stdin.read()
print(json.dumps({
    'hookSpecificOutput': {
        'hookEventName': 'PostToolUse',
        'additionalContext': 'golangci-lint:\n' + data
    }
}))
" <<< "$out"

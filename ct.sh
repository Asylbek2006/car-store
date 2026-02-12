#!/bin/bash
OUT="context.txt"

# Collect project files (only text ones)
find . -type f \
   -not -path '*/.*' \
   -not -path './mockdata/*' \
   -not -name "$(basename "$0")" \
   -not -name "$OUT" \
   -print0 | while IFS= read -r -d $'\0' file; do
   if file "$file" | grep -q 'text'; then
       {
           echo "--- File: $file ---"
           cat "$file"
           echo
       }
   fi
done > "$OUT"

# Append current git diff (if repo)
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
   {
       echo
       echo "--- Git Diff ---"
       git diff
   } >> "$OUT"
fi
#!/bin/bash -eu
cd "$(dirname "${0}")"
VALID_FORMAT='[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*'
TEMP_HELP="$(mktemp -t tc-client-generator-help-text.XXXXXXXXXX)"
TEMP_README="$(mktemp -t tc-client-generator-readme.XXXXXXXXXX)"
TEMP_BINARY="$(mktemp -t tc-client-generator.XXXXXXXXXX)"
go build -a -o "${TEMP_BINARY}" ./cmd/tc-client-generator
"${TEMP_BINARY}" --help > "${TEMP_HELP}"
echo '```' >> "${TEMP_HELP}"
sed -e "
   /^tc-client-generator ${VALID_FORMAT}/,/^\`\`\`\$/!b
   //!d
   /^tc-client-generator ${VALID_FORMAT}/d;r ${TEMP_HELP}
   /^\`\`\`\$/d
" README.md > "${TEMP_README}"
cat "${TEMP_README}" > README.md
rm "${TEMP_BINARY}"
rm "${TEMP_README}"
rm "${TEMP_HELP}"

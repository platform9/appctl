#!/bin/bash
set -eo pipefail
[ ! -z "${BASH_DEBUG}" ] && set -x
PS4='${BASH_SOURCE}.${LINENO}+ '

if ! command -v gimme; then
    echo "ERROR: Please have https://github.com/travis-ci/gimme somewhere on PATH" >&2
    exit 1
fi

eval $(gimme 1.17)

pushd $(dirname $0)
    make --max-load=$(nproc)
popd

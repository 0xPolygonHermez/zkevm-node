#!/bin/sh

set -e

gen() {
    local package=$1

    abigen --abi abi/${package}.abi --pkg=${package} --out=${package}/${package}.go
}

gen uniswap
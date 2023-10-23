#!/bin/sh

set -e

gen() {
    local package=$1

    abigen --bin bin/${package}.bin --abi abi/${package}.abi --pkg=${package} --out=${package}/${package}.go
}

gen polygonzkevm
gen oldpolygonzkevm
gen polygonzkevmbridge
gen oldpolygonzkevmbridge
gen pol
gen polygonzkevmglobalexitroot
gen polygonrollupmanager
gen mockpolygonrollupmanager
gen mockverifier
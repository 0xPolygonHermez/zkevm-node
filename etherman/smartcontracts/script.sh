#!/bin/sh

set -e

gen() {
    local package=$1

    abigen --bin bin/${package}.bin --abi abi/${package}.abi --pkg=${package} --out=${package}/${package}.go
}

genNoBin() {
    local package=$1

    abigen --abi abi/${package}.abi --pkg=${package} --out=${package}/${package}.go
}

gen oldpolygonzkevmglobalexitroot
gen oldpolygonzkevmbridge
gen oldpolygonzkevm
gen etrogpolygonzkevm
gen polygonzkevm
gen polygonzkevmbridge
gen pol
gen polygonzkevmglobalexitroot
gen polygonrollupmanager
gen mockpolygonrollupmanager
gen mockverifier
gen polygondatacommittee
genNoBin dataavailabilityprotocol
gen proxy

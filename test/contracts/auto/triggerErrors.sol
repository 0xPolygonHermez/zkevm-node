// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract triggerErrors {
    uint256 public count = 0;

    // set gasLimit = 50000 & steps = 100
    function outOfGas() public {
        for (uint256 i = 0; i < 100; i++) {
            assembly {
                sstore(0x00, i)
            }
        }
    }

    // set gasLimit = 30000000 & steps = 50000
    function outOfCountersPoseidon() public {
        for (uint256 i = 0; i < 50000; i++) {
            assembly {
                sstore(0x00, i)
            }
        }
    }

    // bytesKeccak = 1000000 & gasLimit = 50000
    function outOfCountersKeccaks() pure public returns (bytes32 test) {
        assembly {
            test := keccak256(0, 1000000)
        }
        return test;
    }

    // set number and gas limit
    // gasLimit = 50000 & iterations = 100000
    function outOfCountersSteps() pure public {
        for (uint i = 0; i < 100000; i++) {
            assembly {
                mstore(0x0, 1234)
            }
        }
    }
}
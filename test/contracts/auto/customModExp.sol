// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <0.9.0;

contract customModExp {
    bytes32 hashResult;
    address retEcrecover;
    bytes dataResult;
    uint256 dataRes;

    bytes32[10] arrayStorage;

    function modExpGeneric(bytes memory input) public {
        bytes32[10] memory output;

        assembly {
            let success := staticcall(gas(), 0x05, add(input, 32), mload(input), output, 0x140)
            sstore(0x00, success)
        }

        for (uint i = 0; i < 10; i++) {
            arrayStorage[i] = output[i];
        }
    }
}
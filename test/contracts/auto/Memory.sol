// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract Memory {
    bool public resStatic = true;

    function testStaticEcrecover() public {
        assembly {
            let resultStatic := staticcall(gas(), 0x01, 0x20, 0x80, 0xa0, 0x20)
            // staticcall(g, a, in, insize, out, outsize) --> input is mem[in...(in + insize)]
            sstore(0, resultStatic)
        }
    }
}
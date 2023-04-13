// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract Create2 {
    function opCreate2(bytes memory bytecode, uint length) public returns(address) {
        address addr;
        assembly {
            addr := create2(0, add(bytecode, 0x20), length, 0x2)
            sstore(0x0, addr)
        }
        return addr;
    }
}
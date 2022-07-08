// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract EmitLog2 {
    event LogA(uint256 indexed a);

    function emitLogs() public {
        emit LogA(1);
    }
}
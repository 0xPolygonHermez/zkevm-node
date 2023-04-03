// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract EmitLog2 {
    event Log();
    event LogA(uint256 indexed a);
    event LogABCD(uint256 indexed a, uint256 indexed b, uint256 indexed c, uint256 d);

    function emitLogs() public {
        assembly {
            log0(0, 32)
        }
        emit Log();
        emit LogA(1);
        emit LogABCD(1, 2, 3, 4);
    }
}
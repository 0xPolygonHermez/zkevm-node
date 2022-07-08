// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract EmitLog {
    event Log();
    event LogA(uint256 indexed a);
    event LogAB(uint256 indexed a, uint256 indexed b);
    event LogABC(uint256 indexed a, uint256 indexed b, uint256 indexed c);
    event LogABCD(uint256 indexed a, uint256 indexed b, uint256 indexed c, uint256 d);

    function emitLogs() public {
        emit LogA(1);
    }
}
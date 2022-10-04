// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ConditionalLoop {
    function ExecuteLoop(uint256 times) external pure returns(uint256) {
        require(times>0, 'times need to be bigger than 0');

        uint256 executed = 0;

        for (uint256 i = 1; i <= times; i++) {
            executed++;
        }

        return executed;
    }
}

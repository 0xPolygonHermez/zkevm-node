// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.4;

contract Revert2 {
    function generateError() public {
        revert("Today is not juernes");
    }
}
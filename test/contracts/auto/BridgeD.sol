// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract BridgeD {
    receive() external payable {}
    
    address sender;
    uint256 value;

    function exec() public payable {
        sender = msg.sender;
        value = msg.value;
    }
}
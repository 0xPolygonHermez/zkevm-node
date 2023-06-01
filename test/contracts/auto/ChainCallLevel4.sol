// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract ChainCallLevel4 {
    address sender;
    uint256 value;
    
    function exec() public payable {
        sender = msg.sender;
        value = msg.value;
    }

    function get() public pure returns (string memory t) {
        return "ahoy";
    }
}

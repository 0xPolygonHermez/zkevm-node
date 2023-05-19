// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract Called {
    uint256 num;
    address sender;
    uint256 value;

    function setVars(uint256 _num) public payable {
        num = _num;
        sender = msg.sender;
        value = msg.value;
    }

    function getVars() public view returns (uint256, address, uint256) {
        return (num, sender, value);
    }
}

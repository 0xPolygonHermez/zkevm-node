// SPDX-License-Identifier: MIT

pragma solidity >=0.7.0 <0.9.0;

// NOTE: Deploy this contract first
contract DelegateCallCalled {
    // NOTE: storage layout must be the same as contract A
    uint public num;
    address public sender;
    uint public value;

    function setVars(uint _num) public payable {
        num = _num;
        sender = msg.sender;
        value = msg.value;
    }
}

// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity ^0.8.0;

contract Delegator {
    // storage of the caller will be copied to the callee.
    address public expectedSender;

    constructor(address _expectedSender) {
        expectedSender = _expectedSender;
    }

    function call(address target) public {
        (bool success,) = target.delegatecall(abi.encodeWithSignature("entrypoint()"));

        require(success, 'expectedSender != msg.sender');
    }
}

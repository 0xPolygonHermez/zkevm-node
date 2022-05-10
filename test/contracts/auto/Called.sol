// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity ^0.8.0;

contract Called {
    // storage from the caller will be copied to the callee.
    address public expectedSender;

    function entrypoint() external view {
        require(expectedSender == msg.sender, 'expectedSender != msg.sender');
    }
}

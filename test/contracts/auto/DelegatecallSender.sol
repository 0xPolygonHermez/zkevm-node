// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity ^0.8.0;

contract DelegatecallSender {
    // storage of the caller will be copied to the callee.
    address public expectedSender;

    function call(address target) public {
        expectedSender = msg.sender;

        (bool success, bytes memory result) = target.delegatecall(abi.encodeWithSignature("entrypoint()"));

        if (!success) {
            // Next 5 lines from https://ethereum.stackexchange.com/a/83577
            if (result.length < 68) revert();
            assembly {
                result := add(result, 0x04)
            }
            revert(abi.decode(result, (string)));
        }
        require(success, 'delegated call failed');
    }
}

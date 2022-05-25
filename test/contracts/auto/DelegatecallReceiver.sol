// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity ^0.8.0;

contract DelegatecallReceiver {
    // storage from the caller will be copied to the callee.
    address public expectedSender;

    bytes16 private constant _HEX_SYMBOLS = "0123456789abcdef";

    function entrypoint() external view {
        require(expectedSender == msg.sender, string(abi.encodePacked('expectedSender ', toHexString(uint160(expectedSender), 20), ' actual sender ', toHexString(uint160(msg.sender), 20))));
    }

    function toHexString(uint256 value, uint256 length) internal pure returns (string memory) {
        bytes memory buffer = new bytes(2 * length + 2);
        buffer[0] = "0";
        buffer[1] = "x";
        for (uint256 i = 2 * length + 1; i > 1; --i) {
            buffer[i] = _HEX_SYMBOLS[value & 0xf];
            value >>= 4;
        }
        require(value == 0, "Strings: hex length insufficient");
        return string(buffer);
    }

}

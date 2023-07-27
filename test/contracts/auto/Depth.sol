// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract Depth {
    uint test = 0;
    bytes32 constant auxReturn = 0x6aecbc3300000000000000000000000000000000000000000000000000000000;

    function start(address addr, uint256 gasForwarded) public {
        test = this.secondCall{gas: gasForwarded}(addr);
    }

    function secondCall(address addr) external returns (uint256) {
        uint256 success;
        assembly {
            mstore(0x80, auxReturn)
            success := staticcall(gas(), addr, 0x80, 0x04, 0x80, 0x20)
        }
    }
}
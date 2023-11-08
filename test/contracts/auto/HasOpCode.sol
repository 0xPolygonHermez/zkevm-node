// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

contract HasOpCode {
    uint256 gasPrice = 0;
    uint256 balance = 0;

    function opGasPrice() public {
        uint256 tmp;
        assembly {
            tmp := gasprice()
        }
        gasPrice = tmp;
    }

    function opBalance() public {
        address a = msg.sender;
        uint256 tmp;
        assembly {
            tmp := balance(a)
        }
        balance = tmp;
    }
}
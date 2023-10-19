// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ConstructorMap {
    mapping(uint => uint) public numbers;

    constructor() {
        uint i = 0;
        for (i = 0; i < 100; i++) {
            numbers[i] = i;
        }
    }
}

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract StorageOnDeploy {

    uint256 number;
    
    constructor() {
        number = 1234;
    }

    function retrieve() public view returns (uint256){
        return number;
    }
}

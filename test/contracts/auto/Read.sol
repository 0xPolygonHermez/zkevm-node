// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract Read {

    uint256 value = 1;

    function publicRead() public view returns (uint256){
        return value;
    }

    function publicReadWParams(uint256 p) public view returns (uint256){
        return value + p;
    }
    
    function externalRead() external view returns (uint256){
        return value;
    }    

    function externalReadWParams(uint256 p) external view returns (uint256){
        return value + p;
    }
}

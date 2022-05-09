// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract FailureTest {
    uint256 number;
    event numberChanged(uint256 from, uint256 to);

    function store(uint256 num) public {
        uint256 oldNum = number;
        number = num;
        emit numberChanged(oldNum, num);
    }

    function storeAndFail(uint256 num) public {
        store(num);
        require(true == false, "this method always fails");
    }

    function getNumber() public view returns (uint256){
        return number;
    }
}
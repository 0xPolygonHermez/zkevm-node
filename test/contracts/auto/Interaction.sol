// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

interface ICounter {
    function count() external view returns (uint);
    function increment() external;
}

contract Interaction {
    address counterAddr;

    function setCounterAddr(address _counter) public payable {
       counterAddr = _counter;
    }

    function getCount() external view returns (uint) {
        return ICounter(counterAddr).count();
    }
}

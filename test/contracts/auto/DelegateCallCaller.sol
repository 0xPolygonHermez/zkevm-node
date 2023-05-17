// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <0.9.0;

contract DelegateCallCaller {
    uint public num;
    address public sender;
    uint public value;

    function setVars(address _contract, uint _num) public payable returns (bool, bytes memory) {
        // A's storage is set, B is not modified.
        (bool success, bytes memory data) = _contract.delegatecall(
            abi.encodeWithSignature("setVars(uint256)", _num)
        );
        return (success, data);
    }
}
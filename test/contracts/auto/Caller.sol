// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract Caller {
    function execCall(address _contract, uint _num) public payable {
        bool ok;
        (ok, ) = _contract.call(
            abi.encodeWithSignature("setVars(uint256)", _num)
        );
        require(ok, "failed to perform call");

        (ok, ) = _contract.delegatecall(
            abi.encodeWithSignature("setVars(uint256)", _num)
        );
        require(ok, "failed to perform delegate call");
        
        bytes memory result;
        (ok, result) = _contract.staticcall(
            abi.encodeWithSignature("getVars()")
        );
        require(ok, "failed to perform static call");

        uint256 num;
        address sender;
        uint256 value;

        (num, sender, value) = abi.decode(result, (uint256, address, uint256));
    }
}

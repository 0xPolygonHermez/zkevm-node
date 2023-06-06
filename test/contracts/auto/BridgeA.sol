// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract BridgeA {
    receive() external payable {}
    
    function exec(address bridgeB, address bridgeC, address bridgeD, address acc) public payable {
        bool ok;
        (ok,) = bridgeB.delegatecall(abi.encodeWithSignature("exec(address,address,address)", bridgeC, bridgeD, acc));
        require(ok, "failed to perform delegate call to bridge B");
    }
}
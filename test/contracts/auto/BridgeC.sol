// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract BridgeC {
    receive() external payable {}
    
    function exec(address bridgeD) public payable {
        bool ok;
        (ok,) = bridgeD.delegatecall(abi.encodeWithSignature("exec()"));
        require(ok, "failed to perform delegate call to bridge D");
    }
}

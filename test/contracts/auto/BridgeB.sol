// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract BridgeB {
    receive() external payable {}
    
    function exec(address bridgeC, address bridgeD, address acc) public payable {
        bool ok;
        (ok,) = bridgeC.call(abi.encodeWithSignature("exec(address)", bridgeD));
        require(ok, "failed to perform call to bridge C");

        (ok,) = acc.call{value:msg.value}("");
        require(ok, "failed to perform call to acc");
    }
}

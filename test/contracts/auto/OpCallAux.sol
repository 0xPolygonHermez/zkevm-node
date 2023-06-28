// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <0.9.0;

contract OpCallAux {
    uint256 auxVal = 1;

    function auxReturn() external returns (uint256) {
        return 0x123456689;
    }

    function auxUpdate() external returns (uint256) {
        assembly {
            sstore(0x0, 0x12121212121212121212)
        }
        return 0x123456689;
    }

    function opDelegateSelfBalance() external payable returns (uint256) {
        auxVal = address(this).balance;
    }

    function opCallSelfBalance() external payable returns (uint256) {
        auxVal = address(this).balance;
    }

    // function opDelegateCallSelfBalance(
    //     address addrCall
    // ) external payable returns (uint256) {
    //     addrCall.call{value: msg.value}(
    //         abi.encodeWithSignature("opDelegateCallSelfBalance()")
    //     );
    //     assembly {
    //         let val := address()
    //         sstore(0x3, val)
    //         let val2 := codesize()
    //         sstore(0x4, val2)
    //         let val3 := msize()
    //         codecopy(val3, 0x0, val2)
    //         let val4 := mload(val3)
    //         sstore(0x5, val3)
    //     }
    // }

    // function opDelegateDelegateSelfBalance(
    //     address addrCall
    // ) external payable returns (uint256) {
    //     addrCall.delegatecall(
    //         abi.encodeWithSignature("opDelegateCallSelfBalance()")
    //     );
    //     assembly {
    //         let val := address()
    //         sstore(0x3, val)
    //         let val2 := codesize()
    //         sstore(0x4, val2)
    //         let val3 := msize()
    //         codecopy(val3, 0x0, val2)
    //         let val4 := mload(val3)
    //         sstore(0x5, val3)
    //     }
    // }

    function addTwo(uint256 a, uint256 b) public returns (uint256) {
        return a + b;
    }

    function opReturnCallSelfBalance(
        address addrCall
    ) external payable returns (uint256) {
        return address(this).balance;
    }

    function auxUpdateValues() external payable returns (uint256) {
        address send = msg.sender;
        uint256 val = msg.value;
        assembly {
            sstore(0x0, 0x12121212121212121212)
            sstore(0x1, send)
            sstore(0x2, val)
        }
        return 0x123456689;
    }

    function auxFail() external {
        require(1 == 0);
    }

    function auxStop() external {
        require(0 == 0);
    }
}
// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract Read {
    struct token {
        string Name;
        uint256 Quantity;
        address Address;
    }

    address public Owner;
    string public OwnerName;
    uint256 public Value = 1;
    mapping(address => token) public Tokens;

    constructor(string memory name) {
        OwnerName = name;
        Owner = msg.sender;
    }

    function publicGetOwnerName() public view returns (string memory) {
        return OwnerName;
    }

    function externalGetOwnerName() external view returns (string memory) {
        return OwnerName;
    }

    function publicAddToken(token memory t) public {
        Tokens[t.Address] = t;
    }

    function externalAddToken(token memory t) external {
        Tokens[t.Address] = t;
    }

    function publicGetToken(address a) public view returns (token memory) {
        return Tokens[a];
    }

    function externalGetToken(address a) external view returns (token memory) {
        return Tokens[a];
    }

    function publicRead() public view returns (uint256) {
        return Value;
    }

    function externalRead() external view returns (uint256) {
        return Value;
    }

    function publicReadWParams(uint256 p) public view returns (uint256) {
        return Value + p;
    }

    function externalReadWParams(uint256 p) external view returns (uint256) {
        return Value + p;
    }
}

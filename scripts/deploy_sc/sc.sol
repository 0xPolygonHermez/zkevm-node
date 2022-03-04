// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

/**
 * @title Event
 * @dev Store & retrieve value in a variable emitting an event
 */
contract SC {
    uint256 value;
    
    event OldValue(uint256 value);
    event NewValue(uint256 value);
    event ValueChanged(uint256 from, uint256 to);
    event Msg(string s);

    constructor() {
        emit Msg("contract created!");
    }

    function store(uint256 newValue) public {
        uint256 oldValue = value;
        value = newValue;
        emit OldValue(oldValue);
        emit NewValue(newValue);
        emit ValueChanged(oldValue, newValue);
    }

    function retrieve() public view returns (uint256){
        return value;
    }
}
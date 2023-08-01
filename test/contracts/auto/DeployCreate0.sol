// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract DeployCreate0 {
   constructor () {
      assembly {
          let addr := create(0,0,0)
      }
   }
}
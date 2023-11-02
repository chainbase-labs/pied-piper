// SPDX-License-Identifier: GPL-3.0

//1. depoly
// forge create  --private-key ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80  --rpc-url http://localhost:8545 --contracts contract/test.sol Storage
//2. set and store
// cast send 0x32cd5ecdA7f2B8633C00A0434DE28Db111E60636  --from 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --unlocked "store(uint256)" 123
// cast call 0x32cd5ecdA7f2B8633C00A0434DE28Db111E60636   "retrieve()uint256"

//3. same address used create2 TODO:https://book.getfoundry.sh/reference/cast/cast-create2?highlight=create2#cast-create2

pragma solidity >0.7.0 < 0.9.0;
/**
 * @title Storage
 * @dev store or retrieve variable value
 */

contract Storage {
    uint256 value;

    function store(uint256 number) public {
        value = number;
    }

    function retrieve() public view returns (uint256) {
        return value;
    }
}

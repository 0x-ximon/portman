// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Address} from "@openzeppelin/contracts/utils/Address.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";

contract PortmanManager {
    using Clones for address;
    using Address for address;

    address private immutable IMPLEMENTATION;
    error ImplementationNotContract();

    constructor(address _implementation) {
        require(_implementation.code.length > 0, ImplementationNotContract());
        IMPLEMENTATION = _implementation;
    }

    function predictAddress(bytes calldata data) public view returns (address) {
        return IMPLEMENTATION.predictDeterministicAddress(keccak256(data), address(this));
    }

    function cloneAndInitialize(bytes calldata data) external returns (address) {
        address predicted = predictAddress(data);
        if (predicted.code.length == 0) {
            IMPLEMENTATION.cloneDeterministic(keccak256(data));
            predicted.functionCall(data);
        }

        return predicted;
    }
}

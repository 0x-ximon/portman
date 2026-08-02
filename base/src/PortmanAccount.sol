// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Account} from "@openzeppelin/contracts/account/Account.sol";
import {IERC1271} from "@openzeppelin/contracts/interfaces/IERC1271.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {MultiSignerERC7913} from "@openzeppelin/contracts/utils/cryptography/signers/MultiSignerERC7913.sol";

contract PortmanAccount is Account, MultiSignerERC7913, Initializable {
    constructor() MultiSignerERC7913(new bytes[](0), 0) {
        _disableInitializers();
    }

    function initialize(bytes[] memory signers, uint64 threshold) external initializer {
        _addSigners(signers);
        _setThreshold(threshold);
    }

    function addSigners(bytes[] memory signers) external onlyEntryPointOrSelf {
        _addSigners(signers);
    }

    function removeSigners(bytes[] memory signers) external onlyEntryPointOrSelf {
        _removeSigners(signers);
    }

    function setThreshold(uint64 threshold) external onlyEntryPointOrSelf {
        _setThreshold(threshold);
    }

    function isValidSignature(bytes32 hash, bytes calldata signature) external view returns (bytes4) {
        return _rawSignatureValidation(hash, signature) ? IERC1271.isValidSignature.selector : bytes4(0xffffffff);
    }
}

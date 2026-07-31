// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Account} from "@openzeppelin/contracts/account/Account.sol";
import {MultiSignerERC7913} from "@openzeppelin/contracts/utils/cryptography/signers/MultiSignerERC7913.sol";

// HELP: You plan on leveraging these signers with a threshold of 1 for a Portman Account
// import {SignerECDSA} from "@openzeppelin/contracts/utils/cryptography/signers/SignerECDSA.sol";
// import {SignerP256} from "@openzeppelin/contracts/utils/cryptography/signers/SignerP256.sol";
// import {SignerZKEmail} from "@openzeppelin/contracts/utils/cryptography/signers/SignerZKEmail.sol";

// TODO: Consider inheriting if ERC721 or ERC1155 would be used.
// import {ERC721Holder} from "@openzeppelin/contracts/token/ERC721/utils/ERC721Holder.sol";
// import {ERC1155Holder} from "@openzeppelin/contracts/token/ERC1155/utils/ERC1155Holder.sol";

contract PortmanAccount is Account, MultiSignerERC7913 {
    constructor(bytes[] memory signers, uint64 threshold) MultiSignerERC7913(signers, threshold) {}
}

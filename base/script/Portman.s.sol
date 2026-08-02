// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {PortmanAccount} from "../src/PortmanAccount.sol";
import {PortmanManager} from "../src/PortmanManager.sol";

contract PortmanScript is Script {
    function run() public {
        vm.startBroadcast();

        address admin = msg.sender;
        bytes[] memory signers = new bytes[](1);
        signers[0] = abi.encodePacked(admin);

        PortmanAccount account = new PortmanAccount(signers);
        PortmanManager manager = new PortmanManager(address(account));

        console.log("PortmanAccount:", address(account));
        console.log("PortmanManager:", address(manager));

        vm.stopBroadcast();
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {PortmanDispenser} from "../src/PortmanDispenser.sol";

contract PortmanDispenserScript is Script {
    function run() public {
        vm.startBroadcast();

        address admin = msg.sender;
        PortmanDispenser dispenser = new PortmanDispenser(admin);
        console.log("PortmanDispenser:", address(dispenser));

        vm.stopBroadcast();
    }
}

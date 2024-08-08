/*
 *  SPDX-License-Identifier: Apache-2.0
 */

'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const { Wallets } = require('fabric-network');
const path = require('path');

async function main() {

    // Main try/catch block
    try {

        // A wallet stores a collection of identities
        const wallet = await Wallets.newFileSystemWallet('wallet');

        // get the user identities of the wallet
        const walletList = await wallet.list();
        console.log('The wallet has these identities:',walletList);

    } catch (error) {
        console.log(`Error adding to wallet. ${error}`);
        console.log(error.stack);
    }
}

main().then(() => {
    console.log('Done');
}).catch((e) => {
    console.log(e);
    console.log(e.stack);
    process.exit(-1);
});

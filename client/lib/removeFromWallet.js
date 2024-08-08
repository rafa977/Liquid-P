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

        var args = process.argv.slice(2);

        if (args.length != 1 ) {
            console.log('Wrong number of parameters found. Expending 6 params:\n#1. the identiy label e.g. "user1_be_tax"\n#2. the org MSP ID e.g "BETAXIOSSDRPOCMSP"\n#3. the path to the user signcerts folder e.g. "../network/organizations/peerOrganizations/betax.ebsi.eu/users/user1_be_tax@betax.ebsi.eu/msp/signcerts"\n#4. the .PEM filename e.g. "cert.pem"\n#5.the path to the user keystore folder e.g. "../network/organizations/peerOrganizations/betax.ebsi.eu/users/user1_be_tax@betax.ebsi.eu/msp/keystore" \n#6. the private key filename of the user e.g.  "33af6f24b7a4d8a317fca3fe7939d45ce801b923b3d918d814d8eff87d4b49df_sk" ');

            return;
        }

        var in_identityLabel = args[0]; //eg "user1_be_tax"

        // A wallet stores a collection of identities
        const wallet = await Wallets.newFileSystemWallet('wallet');

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.get(in_identityLabel);
        if (!userExists) {
            console.log('An identity for the user: '+in_identityLabel+' does not exists in the wallet.');
            console.log(userExists);
        }

        
        await wallet.remove(in_identityLabel);
        console.log('The user identiy is removed from the wallet.');

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

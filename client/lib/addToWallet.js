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

        if (args.length != 6 ) {
            console.log('Wrong number of parameters found. Expending 6 params:\n#1. the identiy label e.g. "user1_be_tax"\n#2. the org MSP ID e.g "BETAXIOSSDRPOCMSP"\n#3. the path to the user signcerts folder e.g. "../network/organizations/peerOrganizations/betax.ebsi.eu/users/user1_be_tax@betax.ebsi.eu/msp/signcerts"\n#4. the .PEM filename e.g. "cert.pem"\n#5.the path to the user keystore folder e.g. "../network/organizations/peerOrganizations/betax.ebsi.eu/users/user1_be_tax@betax.ebsi.eu/msp/keystore" \n#6. the private key filename of the user e.g.  "33af6f24b7a4d8a317fca3fe7939d45ce801b923b3d918d814d8eff87d4b49df_sk" ');

            return;
        }

        var in_identityLabel = args[0]; //eg "user1_be_tax"
        var in_orgMSPID = args[1]; // eg BETAXIOSSDRPOCMSP"

        var in_signcertsPath = args[2]; // eg "../network/organizations/peerOrganizations/betax.ebsi.eu/users/user1_be_tax@betax.ebsi.eu/msp/signcerts"
        var in_certificateFilename = args[3]; // eg "cert.pem"

        var in_keystorePath = args[4]; // eg "../network/organizations/peerOrganizations/betax.ebsi.eu/users/user1_be_tax@betax.ebsi.eu/msp/keystore"
        var in_privateKeyFilename = args[5]; // eg "33af6f24b7a4d8a317fca3fe7939d45ce801b923b3d918d814d8eff87d4b49df_sk"
        console.log(`Processing the following:\n 1. identiy label: ${in_identityLabel}\n 2. org MSP ID : ${in_orgMSPID}\n 3. user signcerts path: ${in_signcertsPath}\n 4. user PEM file name: ${in_certificateFilename}\n 5. user keystore path: ${in_keystorePath}\n6. user private key filename: ${in_privateKeyFilename}` );

        // A wallet stores a collection of identities
        const wallet = await Wallets.newFileSystemWallet('wallet');

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.get(in_identityLabel);
        if (userExists) {
            console.log('An identity for the user: '+in_identityLabel+' already exist in the wallet. It will be updated.');
            console.log(userExists);
        }

        // Identity to credentials to be stored in the wallet
        const signcertsPath = path.resolve(__dirname, in_signcertsPath);
        console.log('signcertsPath='+signcertsPath);
        const certificate = fs.readFileSync(path.join(signcertsPath, '/', in_certificateFilename)).toString();
        console.log('certificate='+certificate);

        const keystorePath = path.resolve(__dirname, in_keystorePath);
        console.log('keystorePath='+keystorePath);
        const privateKey = fs.readFileSync(path.join(keystorePath, '/', in_privateKeyFilename)).toString();
        console.log('privateKey='+privateKey);

        // Load credentials into walletkeystore
        const identityLabel = in_identityLabel;//'Admin@betax.ebsi.eu';


        const identity = {
                    credentials: {
                        certificate,
                        privateKey
                    },
                    mspId: in_orgMSPID,//'BETAXIOSSDRPOCMSP',
                    type: 'X.509'
                }

        await wallet.put(identityLabel, identity);
        console.log('The user identiy is added in the wallet.');

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

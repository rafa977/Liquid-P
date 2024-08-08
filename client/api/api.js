/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Wallets, Gateway, X509Identity } = require('fabric-network');

const fs = require('fs');
const path = require('path');
const yaml = require('js-yaml');
const { v4: uuidv4 } = require('uuid');
const crypto = require('crypto');

const { KJUR, KEYUTIL, X509 } = require('jsrsasign');
// const CryptoJS = require('crypto-js');


const winston = require('winston');
const { format } = winston;
const { combine, label, json, splat, timestamp, printf  } = format;
const FabricCAServices = require('fabric-ca-client');
const { User } = require('fabric-common');
const dotenv = require('dotenv');
const jwt = require('jsonwebtoken');

const plainFormat = printf( ({ level, message, timestamp ,metadata}) => {
  let msg = `${timestamp} [${level}] : ${message} `
  if(metadata) {
	   msg += JSON.stringify(metadata)
  }
  return msg
});

const logger = winston.createLogger({
  format: winston.format.combine(
            format.timestamp(),
            plainFormat
        ),
  transports: [
    new (winston.transports.Console)({
      level: 'debug',
      colorize: true,
      handleExceptions: true,
      label: "api"
    }),
    new (winston.transports.File)({
      level: 'debug',
      colorize: false,
      handleExceptions: true,
      label: "api",
      filename: "./logs/ts-api.log",
      json: false
    })
  ]
});

class Api {
    constructor() {
      logger.info('API: the Api class is intitialized');
    }

    validateToken(token, jwtSecretKey) {
      // Tokens are generally passed in the header of the request
      // Due to security reasons.    
      try {
          const verified = jwt.verify(token, jwtSecretKey);
          if(verified){

            return verified.username
          }else{
              // Access Denied
              throw new Error('Access Denied');
          }
      } catch (error) {
          // Access Denied
          throw new Error('Access Denied');
      }
    }

    async verifyDocument(args) {
      try {
    
          var in_username = args[0];
          var hash = args[1];
          var ownerOfFile = args[2];
          var docId = args[3];
    
          // Create a new file system based wallet for managing identities.
          const walletPath = path.join(process.cwd(), 'wallet');
          const wallet =  await Wallets.newFileSystemWallet(walletPath);
          logger.info(`API/verify document: Wallet path: ${walletPath}`);
    
          // Check to see if we've already enrolled the user.
          const userExists = await wallet.get(in_username);
          if (!userExists) {
              logger.info('API/verify document: An identity for the user: '+in_username+' does not exist in the wallet');
              //return 'An identity for the user: "'+in_username+'" does not exist in the wallet';
              throw new Error('An identity for the user: "'+in_username+'" does not exist in the wallet');
          }
    
          // Create a new gateway for connecting to our peer node.
          const gateway = new Gateway();
          logger.info('API/verify document: new Gateway done');

          // Load connection profile; will be used to locate a gateway
          let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));
          // Set connection options; identity and wallet
          let connectionOptions = {
            identity: in_username,
            wallet: wallet,
            discovery: { enabled:true, asLocalhost: false }
          };

          await gateway.connect(connectionProfile, connectionOptions);
          logger.info('API/verify document: gateway.connect done');

          // Get the network (channel) our contract is deployed to.
          const network = await gateway.getNetwork('channel');
          logger.info('API/verify document: gateway.getNetwork done');

          // Get the contract from the network.
          const contract = network.getContract('liquidsc');
          logger.info('API/verify document: network.getContract done');

          // Evaluate the specified transaction.
          const result = await contract.evaluateTransaction("getDocumentById",docId);
          logger.info(`API/verify document: Transaction has been evaluated, result is: ${result.toString()}`);
          var response = JSON.parse(result.toString())

          console.log(response.hash)
    
          var sigValueBase64 = response.signature
          console.log("SIg Value: " + sigValueBase64)

          const owner = await wallet.get(ownerOfFile);
          if (!owner) {
            logger.info('API/verify document: An identity for the user: '+ownerOfFile+' does not exist in the wallet');
            //return 'An identity for the user: "'+in_username+'" does not exist in the wallet';
            throw new Error('An identity for the user: "'+ownerOfFile+'" does not exist in the wallet');
        }
          var hashToAction = hash
    
          //we get the the hash from the document
          //we get the public key from the user who signed
          // Show info about certificate provided
          const certObj = new X509();
          certObj.readCertPEM(owner.credentials.certificate);
          console.log("Detail of certificate provided")
          console.log("Subject: " + certObj.getSubjectString());
          console.log("Issuer (CA) Subject: " + certObj.getIssuerString());
          console.log("Valid period: " + certObj.getNotBefore() + " to " + certObj.getNotAfter());
          console.log("CA Signature validation: " + certObj.verifySignature(KEYUTIL.getKey(owner.credentials.certificate)));
          console.log("");
    
          var certLoaded = owner.credentials.certificate;
          
          // perform signature checking
          var userPublicKey = KEYUTIL.getKey(certLoaded);
          var recover = new KJUR.crypto.Signature({"alg": "SHA256withECDSA"});
          recover.init(userPublicKey);
          recover.updateHex(hashToAction);
          var getBackSigValueHex = new Buffer(sigValueBase64, 'base64').toString('hex');
          // console.log("Signature verified with certificate provided: " + recover.verify(getBackSigValueHex));

          return recover.verify(getBackSigValueHex);
    
      } catch (error) {
          logger.error(`API/verification failed: ${error}`);
          throw error;
      }
    }  


async getQuery(args) {
  try {

      var in_username = args[0];
      var id = args[1];
      var requestType = args[2];

      logger.info(`API/` + requestType + `: in_username: ${in_username} id: ${id}`);

      // Create a new file system based wallet for managing identities.
      const walletPath = path.join(process.cwd(), 'wallet');
      const wallet =  await Wallets.newFileSystemWallet(walletPath);
      logger.info(`API/` + requestType + `: Wallet path: ${walletPath}`);

      // Check to see if we've already enrolled the user.
      const userExists = await wallet.get(in_username);
      if (!userExists) {
          logger.info('API/` + requestType + `: An identity for the user: '+in_username+' does not exist in the wallet');
          //return 'An identity for the user: "'+in_username+'" does not exist in the wallet';
          throw new Error('An identity for the user: "'+in_username+'" does not exist in the wallet');
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      logger.info('API/verify document: new Gateway done');

      // Load connection profile; will be used to locate a gateway
      let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));
      // Set connection options; identity and wallet
      let connectionOptions = {
        identity: in_username,
        wallet: wallet,
        discovery: { enabled:true, asLocalhost: false }
      };

      await gateway.connect(connectionProfile, connectionOptions);
      logger.info('API/verify document: gateway.connect done');

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('channel');
      logger.info('API/verify document: gateway.getNetwork done');

      // Get the contract from the network.
      const contract = network.getContract('liquidsc');
      logger.info('API/verify document: network.getContract done');

      // Evaluate the specified transaction.
      const result = await contract.evaluateTransaction(requestType,id);
      logger.info(`API/verify document: Transaction has been evaluated, result is: ${result.toString()}`);
      return result;

  } catch (error) {
      logger.error(`API/verify document: Failed to evaluate transaction: ${error}`);
      throw error;
  }
}

async getQueryArray(args) {
  try {

      var requestType = args[0];
      var in_username = args[1];
      var allData = args[2];
  
      logger.info(`API/` + requestType + `: in_username: ${in_username} `);

      // Create a new file system based wallet for managing identities.
      const walletPath = path.join(process.cwd(), 'wallet');
      const wallet =  await Wallets.newFileSystemWallet(walletPath);
      logger.info(`API/` + requestType + `: Wallet path: ${walletPath}`);

      // Check to see if we've already enrolled the user.
      const userExists = await wallet.get(in_username);
      if (!userExists) {
          logger.info('API/` + requestType + `: An identity for the user: '+in_username+' does not exist in the wallet');
          //return 'An identity for the user: "'+in_username+'" does not exist in the wallet';
          throw new Error('An identity for the user: "'+in_username+'" does not exist in the wallet');
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      logger.info('API/verify document: new Gateway done');

      // Load connection profile; will be used to locate a gateway
      let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));
      // Set connection options; identity and wallet
      let connectionOptions = {
        identity: in_username,
        wallet: wallet,
        discovery: { enabled:true, asLocalhost: false }
      };

      await gateway.connect(connectionProfile, connectionOptions);
      logger.info('API/verify document: gateway.connect done');

      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('channel');
      logger.info('API/verify document: gateway.getNetwork done');

      // Get the contract from the network.
      const contract = network.getContract('liquidsc');
      logger.info('API/verify document: network.getContract done');

      // Evaluate the specified transaction.
      const result = await contract.evaluateTransaction(requestType, ...allData);
      logger.info(`API/verify document: Transaction has been evaluated, result is: ${result.toString()}`);
      return result;

  } catch (error) {
      logger.error(`API/verify document: Failed to evaluate transaction: ${error}`);
      throw error;
  }
}

  async postRequestWithJsonData(args) {
      try {
        var requestType = args[0];
        var in_username = args[1];
        var allData = args[2];

        logger.info(`API/` + requestType +`: Json Unparsed: ${allData}`);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet =  await Wallets.newFileSystemWallet(walletPath);
        logger.info(`API/` + requestType +`: Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.get(in_username);
        if (!userExists) {
            logger.info('API/verify document: An identity for the user: '+in_username+' does not exist in the wallet');
            throw new Error('An identity for the user: "'+in_username+'" does not exist in the wallet');
        }

        var sigValueBase64 = ""

        if(requestType == "addDocument" || requestType == "updateDocument"){
          var all = JSON.parse(allData[1])
          console.log(all)

          var hashToAction = all.hash
          console.log(hashToAction)
          var privateKey = ""
          if (userExists && userExists.type === 'X.509') {
              console.log("userPrivateKey 1 : " + userExists.type);
              // access .marker here
              privateKey = userExists.credentials.privateKey
    
              console.log("userPrivateKey 1 : " + privateKey);
          }
    
    
          var sig = new KJUR.crypto.Signature({"alg": "SHA256withECDSA"});
          sig.init(privateKey, "");
          sig.updateHex(hashToAction);
          var sigValueHex = sig.sign();
          console.log("Signature 1 : " + sigValueHex);
          sigValueBase64 = new Buffer(sigValueHex, 'hex').toString('base64');
          console.log("Signature: " + sigValueBase64);
          all.signature = sigValueBase64;
          allData[1] = JSON.stringify(all)
        }

        console.log(allData)

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        logger.info('API/verify document: new Gateway done');

        // Load connection profile; will be used to locate a gateway
        let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));
        // Set connection options; identity and wallet
        let connectionOptions = {
          identity: in_username,
          wallet: wallet,
          discovery: { enabled:true, asLocalhost: false }
        };

        await gateway.connect(connectionProfile, connectionOptions);
        logger.info('API/verify document: gateway.connect done');

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('channel');
        logger.info('API/verify document: gateway.getNetwork done');

        // Get the contract from the network.
        const contract = network.getContract('liquidsc');
        logger.info('API/verify document: network.getContract done');

        // Submit the specified transaction.
        await contract.submitTransaction(requestType, ...allData);
        logger.info('API/verify document: Transaction has been submitted successfully');
        
        // Disconnect from the gateway.
        await gateway.disconnect();

        return 'Transaction has been submitted successfully';

      } catch (error) {
          logger.error(`API/`+ requestType +`: Failed to submit transaction: ${error}`);
          throw error;
      }
  }

  async registerUser(args) {

    try {
      var username = args[0];
      var password = args[1];
      var isSuper = args[2];
      var cauth = args[3];
      var caAdmin = args[4];
      var mspId = args[5];
      
      let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));

      // Create a new CA client for interacting with the CA.
      const caURL = connectionProfile.certificateAuthorities[cauth].url;
      const ca = new FabricCAServices(caURL);
      
      ca.tlsCACerts = connectionProfile.certificateAuthorities[cauth].tlsCACerts.pem;

      // Create a new file system based wallet for managing identities.
      const walletPath = path.join(process.cwd(), 'wallet');
      const wallet = await Wallets.newFileSystemWallet(walletPath);
      console.log(`Wallet path: ${walletPath}`);

      // Check to see if we've already enrolled the user.
      const userIdentity = await wallet.get(username);
      if (userIdentity) {
          console.log('An identity for the user '+ username + ' already exists in the wallet');
          throw new Error('An identity for the user: "'+username+'" already exists in the wallet');
      }

      // Check to see if we've already enrolled the admin user.
      const adminIdentity = await wallet.get(caAdmin);
      if (!adminIdentity) {
          console.log('An identity for the admin user "admin" does not exist in the wallet');
          console.log('Run the enrollAdmin.js application before retrying');
          return;
      }

      // build a user object for authenticating with the CA
      const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
      const adminUser = await provider.getUserContext(adminIdentity, caAdmin);
      
      var secret = null;

      if(isSuper == true){
        secret = await ca.register({
            enrollmentID: username,
            enrollmentSecret: password,
            role: 'client'
            ,
            attrs: [{
              name:"role",
              value:"superadmin",
              ecert:true
            }]
        }, adminUser);
      }else{
        secret = await ca.register({
            enrollmentID: username,
            enrollmentSecret: password,
            role: 'client'
        }, adminUser);
      }
      
      console.log(`Here we are : ${secret}`);


      const enrollment = await ca.enroll({
          enrollmentID: username,
          enrollmentSecret: secret
      });
      const x509Identity = {
          credentials: {
              certificate: enrollment.certificate,
              privateKey: enrollment.key.toBytes(),
          },
          mspId: mspId,
          type: 'X.509',
      };
      await wallet.put(username, x509Identity);
      console.log('Successfully registered and enrolled ' + username + ' and imported it into the wallet');

      return 'User created successfully and enrolled.';

  } catch (error) {
      console.error(`Failed to register user ` + username + `: ${error}`);
      throw error;
  }
}  

async updateUserAttribute(args) {

  try {
    var caller = args[0];
    var username = args[1];
    var password = args[2];
    var isSuper = args[3];
    var cauth = args[4];
    var caAdmin = args[5];
    var mspId = args[6];
    
    let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));

    // Create a new CA client for interacting with the CA.
    const caURL = connectionProfile.certificateAuthorities[cauth].url;
    const ca = new FabricCAServices(caURL);
    
    ca.tlsCACerts = connectionProfile.certificateAuthorities[cauth].tlsCACerts.pem;

    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userIdentity = await wallet.get(username);
    if (!userIdentity) {
        console.log('An identity for the user '+ username + ' does not exists in the wallet');
        throw new Error('An identity for the user: "'+username+'" does not exists in the wallet');
    }

    // Check to see if we've already enrolled the admin user.
    const adminIdentity = await wallet.get(caller);
    if (!adminIdentity) {
        console.log('An identity for the admin user "admin" does not exist in the wallet');
        console.log('Run the enrollAdmin.js application before retrying');
        return;
    }

    // build a user object for authenticating with the CA
    const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
    const adminUser = await provider.getUserContext(adminIdentity, caller);
    console.log("user attributes 111: ",adminUser)

          
    const identityService = ca.newIdentityService();
    identityService.client._caName = "ca_liquid";
    
    console.log("user attributeadas s: ", identityService)

    const retrieveIdentity = await identityService.getOne(username, adminUser)

    console.log("user attributes: ",retrieveIdentity.result.attrs)


    //var enrollment = null;

    // if(isSuper == true){
    //   enrollment = await ca.enroll({
    //     enrollmentID: username,
    //     enrollmentSecret: password
    //     ,
    //       attrs: [{
    //         name:"role",
    //         value:"superadmin",
    //         ecert:true
    //       }]
    //   });
    // }else{
    //   enrollment = await ca.enroll({
    //     enrollmentID: username,
    //     enrollmentSecret: password
    //   })
    // }
    

    // const x509Identity = {
    //     credentials: {
    //         certificate: enrollment.certificate,
    //         privateKey: enrollment.key.toBytes(),
    //     },
    //     mspId: mspId,
    //     type: 'X.509',
    // };
    // await wallet.put(username, x509Identity);
    console.log('Successfully registered and enrolled ' + username + ' and imported it into the wallet');

    return 'User created successfully and enrolled.';

} catch (error) {
    console.error(`Failed to register user ` + username + `: ${error}`);
    throw error;
}
}  


async removeUser(username, ca_definition){
  // Main try/catch block
  try {

    let admin_ca = "admin_" + ca_definition;
    let caName = "ca_" + ca_definition;
    
    let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));

    // Create a new CA client for interacting with the CA.
    const caURL = connectionProfile.certificateAuthorities[caName].url;
    const ca = new FabricCAServices(caURL);
    
    ca.tlsCACerts = connectionProfile.certificateAuthorities[caName].tlsCACerts.pem;

    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userIdentity = await wallet.get(username);
    if (!userIdentity) {
        console.log('An identity for the user '+ username + ' does not exists in the wallet');
        throw new Error('An identity for the user: "'+username+'" already exists in the wallet');
    }

    // Check to see if we've already enrolled the admin user.
    const adminIdentity = await wallet.get(admin_ca);
    if (!adminIdentity) {
        console.log('An identity for the admin user does not exist in the wallet');
        //console.log('Run the enrollAdmin.js application before retrying');
        return;
    }

    const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
    const adminUser = await provider.getUserContext(adminIdentity, admin_ca);

    const identityService = ca.newIdentityService();
    identityService.client._caName = caName;

    // console.log(identityService);

    identityService.delete(username, adminUser);

    await wallet.remove(username);
    console.log('The user identiy is removed from the wallet.');

    // get the user identities of the wallet
    const walletList = await wallet.list();
    console.log('The wallet has these identities:',walletList);

  } catch (error) {
      console.log(`Error adding to wallet. ${error}`);
      console.log(error.stack);
  }
}


async revokerUser(username){
  // Main try/catch block
  try {
    let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));

    // Create a new CA client for interacting with the CA.
    const caURL = connectionProfile.certificateAuthorities['ca_liquid'].url;
    const ca = new FabricCAServices(caURL);
    
    ca.tlsCACerts = connectionProfile.certificateAuthorities['ca_liquid'].tlsCACerts.pem;

    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userIdentity = await wallet.get(username);
    if (!userIdentity) {
        console.log('An identity for the user '+ username + ' does not exists in the wallet');
        throw new Error('An identity for the user: "'+username+'" already exists in the wallet');
    }

    // Check to see if we've already enrolled the admin user.
    const adminIdentity = await wallet.get('admin');
    if (!adminIdentity) {
        console.log('An identity for the admin user "admin" does not exist in the wallet');
        console.log('Run the enrollAdmin.js application before retrying');
        return;
    }

    // build a user object for authenticating with the CA
    const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
    const adminUser = await provider.getUserContext(adminIdentity, 'admin');

    var certificate = userIdentity.certificate;
    console.log(certificate)
    

    const secret = null;
    // Register the user, enroll the user, and import the new identity into the wallet.
      secret = await ca.revoke({
        enrollmentID: username,
        aki: "10:6E:3C:D1:D0:53:6C:02:99:46:E4:90:B7:12:CF:05:25:63:0E:44",
        serial: "2c:2f:1a:35:99:84:d1:a2:28:b6:1d:54:1a:21:61:31:95:1d:61:22",
        reason: '8',
        gencrl : false
    }, adminUser);

  } catch (error) {
      console.log(`Error adding to wallet. ${error}`);
      console.log(error.stack);
  }
}


async removeFromWallet(in_identityLabel){
  // Main try/catch block
  try {

    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);
    
  
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

async enrollAdmin(args){
  try {

    var cauth = args[0];
    var adminUsername = args[1];
    var mspId = args[2];

    // load the network configuration
    let connectionProfile = yaml.safeLoad(fs.readFileSync('gateway/connection.yaml', 'utf8'));

    // // load the network configuration
    // const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
    // const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

    // Create a new CA client for interacting with the CA.
    const caURL = connectionProfile.certificateAuthorities[cauth].url;
    // Create a new CA client for interacting with the CA.
    const caInfo = connectionProfile.certificateAuthorities[cauth];
    const caTLSCACerts = caInfo.tlsCACerts.pem;
    const ca = new FabricCAServices(caURL, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);
    
    // Enroll the admin user, and import the new identity into the wallet.
    const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
    const x509Identity = {
        credentials: {
            certificate: enrollment.certificate,
            privateKey: enrollment.key.toBytes(),
        },
        mspId: mspId,
        type: 'X.509',
    };
    await wallet.put(adminUsername, x509Identity);
    console.log('Successfully enrolled admin user ' +adminUsername+' and imported it into the wallet');

    return 'Successfully enrolled admin user.';

} catch (error) {
    console.error(`Failed to enroll admin user ` +adminUsername+`: ${error}`);
    throw error;
}
}

}
module.exports = Api;

# Introduction

The Gateway Application will be deployed as docker container.

The docker image of the Gateway Application will be built based on the provided source files. The deployment team's task list is the following:
- Receive the Gateway Application source files and build the docker image of Gateway Application.
- Change/adapt the connection profile file in order to reflect the target network.
- Copy the connection profile to all Gateway Application nodes.
- Launch the docker container by setting the appropriate values.
- Load the user certificates into Gateway Application wallet for every PoC organization installation.


### Directory structure of source files
The Gateway Application source files should be deployed in the `${PWD}/application` target folder.

The complete directory structure and its contents for the `application` folder is the following:

```
liquid
  ├── network
  ├── chaincode  
  ├── client
  |   ├── addUsersToWallet.sh
  |   ├── api
  |   |   └── api.js
  |   ├── docker
  |   |   ├── .dockerignore
  |   |   ├── Dockerfile
  |   |   └── README.md
  |   ├── gateway
  |   |   └── connection.yaml    
  |   ├── lib
  |   |   ├── addToWallet.js    
  |   |   └── listWallet.js
  |   |   └── removeFromWallet.js      
  |   ├── wallet
  |   ├── logs
  |   |   ├── liquid-api.log    
  |   |   └── liquid-http.log       
  |   |   └── liquid-http_user.log       
  |   ├── https.js
  |   ├── https_user.js
  |   ├── package.json
  |   └── README.md
  └── .env    
```

The `logs` folder will be created and populated automatically by the application (within the docker container internal filesystem). It is shown above in order to provide the full file structure.

## Deploy the source files to the target machine

`mkdir client`

Copy the source files manually to the target folder: `${PWD}/client` as below:

```
https.js        --> ${PWD}/client

https_user.js    --> ${PWD}/client

addUsersToWallet.sh --> ${PWD}/client

package.json    --> ${PWD}/client

api/api.js      --> ${PWD}/client/api

gateway/*.*     --> ${PWD}/client/gateway

lib/*.*         --> ${PWD}/client/lib

docker/*.*         --> ${PWD}/client/docker
```

## Update the "connection profile" to reflect your network

Update the paths and addresses in the connection profile file `client/gateway/connection.yaml` to match the PoC blockchain network on EBSI.

**Attention:** Please be sure that the paths of `tlsCACerts` on the above connection profile start with `.network/` since the container mounts this as volume.

 The `connection.yaml` file has to be copied to all the nodes where the Gateway Application will be installed and run. This file will be mounted on the Gateway Application container.


## Launch Gateway Application as docker container with docker-compose (Build included)

Go to the approperiate directory
`cd /etc/hyperledger/liquid/client`

And run
`docker-compsoe -f docker-compose-app.yaml up -d --build`


# To run the gateway manually
## Build Gateway Application as docker image
Build the docker image from the following location:

`cd /etc/hyperledger/liquid/client`

`docker build -f ./docker/Dockerfile -t liquid-app .`

Check loaded image:
`docker image ls`

Export docker image as .tar:
`docker save -o ./docker/liquid-app.tar  liquid-gateway-app`



## Launch Gateway Application as docker container (Manual Operation)

```
docker run --name liquid-app -p 8081:3000  -v /etc/hyperledger/liquid/network:/usr/src/app/network  -v /etc/hyperledger/liquid/client/gateway:/usr/src/app/gateway  -v /etc/hyperledger/liquid/client/wallet/organizationts:/usr/src/app/wallet  --network=net_liquid -d liquid-app
```

Check that the container is started and get the container ID:

`docker ps`

## Add Admin user from Liquid organization to the wallet (automatic procedure)

`./addUsersToWallet.js`


## To enter the container

`docker exec -it $(docker ps --filter name=liquid-app --format "{{.ID}}") /bin/bash`

Please be sure that you are at : `/usr/src/app` within the container CLI.

Run the following command to list all user identities in the wallet:

`node ./lib/listWallet.js`


## To add users in the Wallet
`docker exec -it $(docker ps --filter name=liquid-app --format "{{.ID}}") /bin/bash`

Get the private key from the network
`ls network/organizations/peerOrganizations/liquid/users/Admin@liquid/msp/keystore`

Replace the private key with the variable that ends with _sk
Replace the corresponding names (Username & MSP)

`node ./lib/addToWallet.js Admin LiquidMSP "../network/organizations/peerOrganizations/liquid/users/Admin@liquid/msp/signcerts" "cert.pem" "../network/organizations/peerOrganizations/liquid/users/Admin@liquid/msp/keystore" "1a43d76a464276452049967b78ed4177d4363d9f35d6421aee49e8399b0807ff_sk"`


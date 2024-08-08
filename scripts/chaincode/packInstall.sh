#!/bin/bash

orderer=$(cat ./network/docker/allOrgs.yaml | shyaml get-value Orderer.Name) > /dev/null

export PATH=/etc/hyperledger/bin:${PWD}:$PATH
export ORDERER_CA=${PWD}/network/organizations/ordererOrganizations/$orderer/orderers/orderer.$orderer/msp/tlscacerts/tlsca.$orderer-cert.pem
export FABRIC_CFG_PATH=/etc/hyperledger/config

rm -rf ./chaincode/liquidsc.tar.gz
peer lifecycle chaincode package ./chaincode/liquidsc.tar.gz --path ./chaincode --lang golang  --label liquidsc

peerOrganizationData=$(./parse_all_peers.sh)

IFS=' ' read -r -a peerOrgs <<< "$peerOrganizationData"
identifier=""

for value in "${peerOrgs[@]}"
do

IFS=',' read -r -a peerInfo <<< "$value"

echo "Peers: " $value

echo '========== Installing chaincode to all peers of each organization =========='
echo '    '

echo '========== Installing chaincode to '${peerInfo[1]}' organization =========='
echo '    '

mspid=${peerInfo[0]}
name=$(echo ${peerInfo[1]} | tr '[:upper:]' '[:lower:]')
port=${peerInfo[2]}

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=$mspid
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/$name/peers/peer0.$name/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/$name/users/Admin@$name/msp
export CORE_PEER_ADDRESS=localhost:$port

#peer lifecycle chaincode install ./chaincode/liquidsc.tar.gz

identifierOP=$(peer lifecycle chaincode install ./chaincode/liquidsc.tar.gz)

IFS=' ' read -r -a identity <<< "$identifierOP"
    for value in "${identity[@]}"
    do
        if [[ $value == *"liquidsc"* ]]; then
            $identifier=$value
            echo $identifier
        fi
    done
done

export CHAINCODE_LABEL=$label;
echo $identifier;


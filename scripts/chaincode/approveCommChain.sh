#!/bin/bash

while getopts a:k:c: flag
do
    case "${flag}" in
        a) chaincodeId=${OPTARG};;
        k) channel=${OPTARG};;
        c) chaincodName=${OPTARG};;
    esac
done

committString=""

orderer=$(cat ./network/docker/allOrgs.yaml | shyaml get-value Orderer.Name) > /dev/null

export PATH=/etc/hyperledger/bin:${PWD}:$PATH
export ORDERER_CA=${PWD}/network/organizations/ordererOrganizations/$orderer/orderers/orderer.$orderer/msp/tlscacerts/tlsca.$orderer-cert.pem
export FABRIC_CFG_PATH=/etc/hyperledger/config


seq=$(peer lifecycle chaincode querycommitted -C $channel -n $chaincodName)
seq=$(echo $seq | cut -d ',' -f 2)
seq=$(echo $seq | cut -d ':' -f 2 )
seq=$((seq+1))
echo $seq

peerOrganizationData=$(./parse_approve_peers.sh)

IFS=' ' read -r -a peerOrgs <<< "$peerOrganizationData"

for value in "${peerOrgs[@]}"
do

IFS=',' read -r -a peerInfo <<< "$value"

echo '========== Approving chaincode to '${peerInfo[1]}' organization =========='
echo '    '

mspid=${peerInfo[0]}
name=$(echo ${peerInfo[1]} | tr '[:upper:]' '[:lower:]')
port=${peerInfo[2]}

peerAddress="--peerAddresses localhost:"$port
tlsCert="--tlsRootCertFiles ${PWD}/network/organizations/peerOrganizations/"$name"/peers/peer0."$name"/tls/ca.crt"
onePeer="$peerAddress $tlsCert"
committString="$committString $onePeer";

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=$mspid
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/$name/peers/peer0.$name/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/$name/users/Admin@$name/msp
export CORE_PEER_ADDRESS=localhost:$port

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.$orderer --tls --cafile $ORDERER_CA --channelID $channel --name $chaincodName --version 1.0 --package-id  $chaincodeId --sequence $seq --signature-policy "AND('ApplicantMSP.member', 'LiquidMSP.member', 'FinancerMSP.member' , 'AltfinancerMSP.member')"
done


peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.authorizer --tls --cafile $ORDERER_CA --channelID $channel --name $chaincodName $committString --version 1.0 --sequence $seq --signature-policy "AND('ApplicantMSP.member', 'LiquidMSP.member', 'FinancerMSP.member' , 'AltfinancerMSP.member')"


echo ' Probably all went well. Enjoy your absolutely perfect written chaincode with no errors :)'

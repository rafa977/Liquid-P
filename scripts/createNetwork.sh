#### Set "ORDERER_CA" this is for our convenience

export ORDERER_CA=${PWD}/network/organizations/ordererOrganizations/authorizer/orderers/orderer.authorizer/msp/tlscacerts/tlsca.authorizer-cert.pem


#### Set Liquid peer0 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="LiquidMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/liquid/peers/peer0.liquid/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/liquid/users/Admin@liquid/msp
export CORE_PEER_ADDRESS=localhost:6051


#### Create the channel

peer channel create  -o localhost:7050  -c channel --ordererTLSHostnameOverride orderer.authorizer  -f ${PWD}/network/channel-artifacts/liquid.tx --outputBlock ./network/channel-artifacts/liquid.block --tls true --cafile ${PWD}/network/organizations/ordererOrganizations/authorizer/orderers/orderer.authorizer/msp/tlscacerts/tlsca.authorizer-cert.pem


### Liquid peer0 join the channel
###Use the ENV variabes from above for BTC organization

#### Join the channel Liquid peer0

peer channel join -b ./network/channel-artifacts/liquid.block


#### Set Liquid peer1 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="LiquidMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/liquid/peers/peer1.liquid/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/liquid/users/Admin@liquid/msp
export CORE_PEER_ADDRESS=localhost:6151


#### Join the channel Liquid peer1

peer channel join -b ./network/channel-artifacts/liquid.block





#### Set Applicant peer0 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="ApplicantMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/applicant/peers/peer0.applicant/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/applicant/users/Admin@applicant/msp
export CORE_PEER_ADDRESS=localhost:7051


#### Join the channel Applicant peer0

peer channel join -b ./network/channel-artifacts/liquid.block


#### Set Applicant peer1 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="ApplicantMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/applicant/peers/peer1.applicant/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/applicant/users/Admin@applicant/msp
export CORE_PEER_ADDRESS=localhost:7151


#### Join the channel Applicant peer1

peer channel join -b ./network/channel-artifacts/liquid.block






#### Set Financer peer0 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="FinancerMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/financer/peers/peer0.financer/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/financer/users/Admin@financer/msp
export CORE_PEER_ADDRESS=localhost:8051


#### Join the channel Financer peer0

peer channel join -b ./network/channel-artifacts/liquid.block


#### Set Financer peer1 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="FinancerMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/financer/peers/peer1.financer/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/financer/users/Admin@financer/msp
export CORE_PEER_ADDRESS=localhost:8151


#### Join the channel Financer peer1

peer channel join -b ./network/channel-artifacts/liquid.block






#### Set AltFinancer peer0 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="AltfinancerMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/altfinancer/peers/peer0.altfinancer/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/altfinancer/users/Admin@altfinancer/msp
export CORE_PEER_ADDRESS=localhost:8251


#### Join the channel Altfinancer peer0

peer channel join -b ./network/channel-artifacts/liquid.block


#### Set Altfinancer peer1 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="AltfinancerMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/altfinancer/peers/peer1.altfinancer/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/altfinancer/users/Admin@altfinancer/msp
export CORE_PEER_ADDRESS=localhost:8351


#### Join the channel Altfinancer peer1

peer channel join -b ./network/channel-artifacts/liquid.block





#### Set Auditor peer0 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="AuditorMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/auditor/peers/peer0.auditor/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/auditor/users/Admin@auditor/msp
export CORE_PEER_ADDRESS=localhost:9051


#### Join the channel Auditor peer0

peer channel join -b ./network/channel-artifacts/liquid.block


#### Set Auditor peer1 variabes

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="AuditorMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/auditor/peers/peer1.auditor/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/auditor/users/Admin@auditor/msp
export CORE_PEER_ADDRESS=localhost:9151


#### Join the channel Auditor peer1

peer channel join -b ./network/channel-artifacts/liquid.block




### UPDATE ANCHOR PEERS 

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="LiquidMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/liquid/peers/peer0.liquid/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/liquid/users/Admin@liquid/msp
export CORE_PEER_ADDRESS=localhost:6051

peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.authorizer -c channel -f ${PWD}/network/channel-artifacts/LiquidAnchors.tx --tls --cafile $ORDERER_CA



export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="ApplicantMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/applicant/peers/peer0.applicant/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/applicant/users/Admin@applicant/msp
export CORE_PEER_ADDRESS=localhost:7051

peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.authorizer -c channel -f ${PWD}/network/channel-artifacts/ApplicantAnchors.tx --tls --cafile $ORDERER_CA



export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="FinancerMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/financer/peers/peer0.financer/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/financer/users/Admin@financer/msp
export CORE_PEER_ADDRESS=localhost:8051

peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.authorizer -c channel -f ${PWD}/network/channel-artifacts/FinancerAnchors.tx --tls --cafile $ORDERER_CA



export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="AltfinancerMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/altfinancer/peers/peer0.altfinancer/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/altfinancer/users/Admin@altfinancer/msp
export CORE_PEER_ADDRESS=localhost:8251

peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.authorizer -c channel -f ${PWD}/network/channel-artifacts/AltfinancerAnchors.tx --tls --cafile $ORDERER_CA




export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="AuditorMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/auditor/peers/peer0.auditor/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/auditor/users/Admin@auditor/msp
export CORE_PEER_ADDRESS=localhost:9051

peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.authorizer -c channel -f ${PWD}/network/channel-artifacts/AuditorAnchors.tx --tls --cafile $ORDERER_CA

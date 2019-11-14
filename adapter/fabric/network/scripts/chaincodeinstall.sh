#!/bin/bash

export FABRIC_CFG_PATH=/etc/hyperledger/fabric
echo $FABRIC_CFG_PATH

CC_NAME=log
CC_PATH=github.com/chaincode
VER=1
CHANNEL_NAME=logchannel
DELAY=5
COUNTER=1
MAX_RETRY=20

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/logisticstransfer.com/orderers/orderer.logisticstransfer.com/msp/tlscacerts/tlsca.logisticstransfer.com-cert.pem

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED to execute current Scenario ==========="
    echo
    exit 1
  fi
}

# peer0.seller Installing chaincode in
echo "========== Installing chaincode in peer0.seller.com to channel $CHANNEL_NAME =========="
sleep 20
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.com/users/Admin@seller.com/msp
CORE_PEER_ADDRESS=peer0.seller.com:7051
CORE_PEER_LOCALMSPID=sellerMSP
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.com/peers/peer0.seller.com/tls/ca.crt
echo "CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS"
peer chaincode install -n $CC_NAME -v $VER -l golang -p $CC_PATH

# peer1.seller Installing chaincode in
echo "========== Installing chaincode in peer1.seller.com to channel $CHANNEL_NAME =========="
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.com/users/Admin@seller.com/msp
CORE_PEER_ADDRESS=peer1.seller.com:8051
CORE_PEER_LOCALMSPID=sellerMSP
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.com/peers/peer1.seller.com/tls/ca.crt
echo "CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS"
peer chaincode install -n $CC_NAME -v $VER -l golang -p $CC_PATH

# peer0.logis Installing chaincode in
echo "========== Installing chaincode in peer0.logis.com to channel $CHANNEL_NAME =========="
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logis.com/users/Admin@logis.com/msp
CORE_PEER_ADDRESS=peer0.logis.com:9051
CORE_PEER_LOCALMSPID=logisMSP
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logis.com/peers/peer0.logis.com/tls/ca.crt
echo "CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS"
peer chaincode install -n $CC_NAME -v $VER -l golang -p $CC_PATH

# peer1.logis Installing chaincode in
echo "========== Installing chaincode in peer1.logis.com to channel $CHANNEL_NAME =========="
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logis.com/users/Admin@logis.com/msp
CORE_PEER_ADDRESS=peer1.logis.com:10051
CORE_PEER_LOCALMSPID=logisMSP
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logis.com/peers/peer1.logis.com/tls/ca.crt
echo "CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS"
peer chaincode install -n $CC_NAME -v $VER -l golang -p $CC_PATH





# peer0.buyer Installing chaincode in
echo "========== Installing chaincode in peer0.buyer.com to channel $CHANNEL_NAME =========="
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.com/users/Admin@buyer.com/msp
CORE_PEER_ADDRESS=peer0.buyer.com:11051
CORE_PEER_LOCALMSPID=buyerMSP
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.com/peers/peer0.buyer.com/tls/ca.crt
echo "CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS"
peer chaincode install -n $CC_NAME -v $VER -l golang -p $CC_PATH

# peer1.buyer Installing chaincode in
echo "========== Installing chaincode in peer1.buyer.com to channel $CHANNEL_NAME =========="
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.com/users/Admin@buyer.com/msp
CORE_PEER_ADDRESS=peer1.buyer.com:12051
CORE_PEER_LOCALMSPID=buyerMSP
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.com/peers/peer1.buyer.com/tls/ca.crt
echo "CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS"
peer chaincode install -n $CC_NAME -v $VER -l golang -p $CC_PATH


echo""
echo "==================================== DONE ======================================"
echo""


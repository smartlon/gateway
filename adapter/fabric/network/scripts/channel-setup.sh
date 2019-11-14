#!/bin/bash

export FABRIC_CFG_PATH=/etc/hyperledger/fabric
echo $FABRIC_CFG_PATH

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

joinChannelWithRetry() {
  PEER=$1
  ORG=$2
  source .config.core.txt

  set -x
  peer channel join -b $CHANNEL_NAME.block >&log.txt
  res=$?
  set +x
  cat log.txt
  if [ $res -ne 0 -a $COUNTER -lt $MAX_RETRY ]; then
    COUNTER=$(expr $COUNTER + 1)
    echo "${PEER}.${ORG}.com failed to join the channel, Retry after $DELAY seconds"
    sleep $DELAY
    joinChannelWithRetry $PEER $ORG
  else
    COUNTER=1
  fi
  verifyResult $res "After $MAX_RETRY attempts, ${PEER}.${ORG}.com has failed to join channel '$CHANNEL_NAME' "
}

# Channel creation
echo "========== Creating channel: "$CHANNEL_NAME" =========="
#sleep 20
peer channel create -o orderer.logisticstransfer.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile $ORDERER_CA

# peer0.seller channel join
echo "========== Joining peer0.seller.com to channel $CHANNEL_NAME =========="
sleep 20
echo "export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.com/users/Admin@seller.com/msp" > .config.core.txt
echo "export CORE_PEER_ADDRESS=peer0.seller.com:7051" >> .config.core.txt
echo "export CORE_PEER_LOCALMSPID=sellerMSP" >> .config.core.txt
echo "export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.com/peers/peer0.seller.com/tls/ca.crt" >> .config.core.txt
joinChannelWithRetry "peer0" "seller" 
peer channel update -o orderer.logisticstransfer.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA

# peer1.seller channel join
echo "========== Joining peer1.seller.com to channel $CHANNEL_NAME =========="
echo "export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.com/users/Admin@seller.com/msp" > .config.core.txt
echo "export CORE_PEER_ADDRESS=peer1.seller.com:8051" >> .config.core.txt
echo "export CORE_PEER_LOCALMSPID=sellerMSP" >> .config.core.txt
echo "export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.com/peers/peer1.seller.com/tls/ca.crt" >> .config.core.txt
joinChannelWithRetry "peer1" "seller"

# peer0.logis channel join
echo "========== Joining peer0.logis.com to channel $CHANNEL_NAME =========="
echo "export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logis.com/users/Admin@logis.com/msp" > .config.core.txt
echo "export CORE_PEER_ADDRESS=peer0.logis.com:9051" >> .config.core.txt
echo "export CORE_PEER_LOCALMSPID=logisMSP" >> .config.core.txt
echo "export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logis.com/peers/peer0.logis.com/tls/ca.crt" >> .config.core.txt
joinChannelWithRetry "peer0" "logis"
peer channel update -o orderer.logisticstransfer.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA

# peer1.logis channel join
echo "========== Joining peer1.logis.com to channel $CHANNEL_NAME =========="
echo "export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logis.com/users/Admin@logis.com/msp" > .config.core.txt
echo "export CORE_PEER_ADDRESS=peer1.logis.com:10051" >> .config.core.txt
echo "export CORE_PEER_LOCALMSPID=logisMSP" >> .config.core.txt
echo "export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logis.com/peers/peer1.logis.com/tls/ca.crt" >> .config.core.txt
joinChannelWithRetry "peer1" "logis"


# peer0.buyer channel join
echo "========== Joining peer0.buyer.com to channel $CHANNEL_NAME =========="
echo "export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.com/users/Admin@buyer.com/msp" > .config.core.txt
echo "export CORE_PEER_ADDRESS=peer0.buyer.com:11051" >> .config.core.txt
echo "export CORE_PEER_LOCALMSPID=buyerMSP" >> .config.core.txt
echo "export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.com/peers/peer0.buyer.com/tls/ca.crt" >> .config.core.txt
joinChannelWithRetry "peer0" "buyer"
peer channel update -o orderer.logisticstransfer.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA

# peer1.buyer channel join
echo "========== Joining peer1.buyer.com to channel $CHANNEL_NAME =========="
echo "export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.com/users/Admin@buyer.com/msp" > .config.core.txt
echo "export CORE_PEER_ADDRESS=peer1.buyer.com:12051" >> .config.core.txt
echo "export CORE_PEER_LOCALMSPID=buyerMSP" >> .config.core.txt
echo "export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.com/peers/peer1.buyer.com/tls/ca.crt" >> .config.core.txt
joinChannelWithRetry "peer1" "buyer"

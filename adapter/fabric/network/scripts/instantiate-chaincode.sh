#!/bin/bash

export FABRIC_CFG_PATH=/etc/hyperledger/fabric


CC_NAME=log
VER=1
CHANNEL_NAME=logchannel

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/logisticstransfer.com/orderers/orderer.logisticstransfer.com/msp/tlscacerts/tlsca.logisticstransfer.com-cert.pem

echo "========== Instantiating chaincode v$VER =========="
peer chaincode instantiate -o orderer.logisticstransfer.com:7050  \
                           --tls $CORE_PEER_TLS_ENABLED     \
                           --cafile $ORDERER_CA             \
                           -C $CHANNEL_NAME                 \
                           -n $CC_NAME                      \
                           -c '{"Args": ["Init"]}'          \
                           -v $VER                          \
			               -l golang                        \
                           -P "AND ('sellerMSP.member','logisMSP.member','buyerMSP.member')"

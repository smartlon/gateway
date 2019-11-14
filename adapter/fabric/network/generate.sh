#!/bin/bash


export IMAGE_TAG=latest
export PATH=${PWD}/../bin:${PWD}:$PATH

echo "************************Generating of Crypto-config certificates***"
cryptogen generate --config=./crypto-config.yaml

echo "************************Generation of channel-artifacts*************"

mkdir channel-artifacts

configtxgen -profile ThreeOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID logchannel
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/logisMSPanchors.tx -channelID  logchannel -asOrg logisMSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/buyerMSPanchors.tx -channelID  logchannel -asOrg buyerMSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/sellerMSPanchors.tx -channelID  logchannel -asOrg sellerMSP



for Org  in seller logis buyer ;
do
cp $Org.yaml.tmp $Org.yaml 
done


ARCH=$(uname -s | grep Darwin)
if [ "$ARCH" == "Darwin" ]; then
  OPTS="-it"
  rm -rf *.yamlt
else
  OPTS="-i"
fi


CURRENT_DIR=$PWD
cd crypto-config/peerOrganizations/seller.com/ca/
PRIV_KEY=$(ls *_sk)
cd "$CURRENT_DIR"
sed $OPTS "s/CA_seller_KEY/${PRIV_KEY}/g" seller.yaml

CURRENT_DIR=$PWD
cd crypto-config/peerOrganizations/logis.com/ca/
PRIV_KEY=$(ls *_sk)
cd "$CURRENT_DIR"
sed $OPTS "s/CA_logis_KEY/${PRIV_KEY}/g" logis.yaml

CURRENT_DIR=$PWD
cd crypto-config/peerOrganizations/buyer.com/ca/
PRIV_KEY=$(ls *_sk)
cd "$CURRENT_DIR"
sed $OPTS "s/CA_buyer_KEY/${PRIV_KEY}/g" buyer.yaml
#!/bin/bash

toLowerCase () {
  local input=$1
  echo $input | tr '[:upper:]' '[:lower:]'
}

L1UpgradeKeyAddress=`toLowerCase $1`

PathToDeployment=$2

getContractAddress () {
  local jsonPath=$1
  local result=`jq .address $jsonPath`
  echo `toLowerCase ${result//\"/}`
}

stripLeadingZeros () {
  local input=$1
  echo "0x${input##0x000000000000000000000000}"
}

callAddressGetter () {
  local contractAddress=$1
  local getFun=$2
  local address=`cast call --flashbots $contractAddress "${getFun}()" ""`
  toLowerCase `stripLeadingZeros ${address}`
}


echo "[PASSED] L1UpgradeKeyAddress $L1UpgradeKeyAddress"

SystemDictatorProxy=`getContractAddress $PathToDeployment/SystemDictatorProxy.json`
echo -n "    Checking SystemDictatorProxy.owner() === L1UpgradeKeyAddress ..."
if [[ xx"$L1UpgradeKeyAddress" == xx`callAddressGetter $SystemDictatorProxy owner` ]]; then
  echo YES
else
  echo NO;
  exit -1
fi

echo -n "    Checking SystemDictatorProxy.admin() === L1UpgradeKeyAddress ..."
if [[ xx"$L1UpgradeKeyAddress" == xx`callAddressGetter $SystemDictatorProxy admin` ]]; then
  echo YES
else
  echo NO;
  exit -1
fi
echo "[PASSED] SystemDictatorProxy $SystemDictatorProxy"

L1ProxyAdmin=`getContractAddress $PathToDeployment/ProxyAdmin.json`
echo -n "    Checking L1ProxyAdmin.owner() === SystemDictatorProxy ..."
if [[ xx"$SystemDictatorProxy" == xx`callAddressGetter $L1ProxyAdmin owner` ]]; then
  echo YES
else
  echo NO;
  # exit -1
fi
echo "[PASSED] L1ProxyAdmin $L1ProxyAdmin"

AddressManager=`getContractAddress $PathToDeployment/Lib_AddressManager.json`
echo -n "    Checking AddressManager.owner() === L1UpgradeKeyAddress ..."
if [[ xx"$L1UpgradeKeyAddress" == xx`callAddressGetter $AddressManager owner` ]]; then
  echo YES
else
  echo NO;
  # exit -1
fi
echo "[PASSED] AddressManager $AddressManager"

L1CrossDomainMessengerProxy=`getContractAddress $PathToDeployment/Proxy__OVM_L1CrossDomainMessenger.json`
echo "[PASSED] L1CrossDomainMessengerProxy $L1CrossDomainMessengerProxy"


L1StandardBridgeProxy=`getContractAddress $PathToDeployment/Proxy__OVM_L1StandardBridge.json`
echo -n "    Checking L1StandardBridgeProxy.getOwner() === L1UpgradeKeyAddress ..."
if [[ xx"$L1UpgradeKeyAddress" == xx`callAddressGetter $L1StandardBridgeProxy getOwner` ]]; then
  echo YES
else
  echo NO;
  # exit -1
fi
echo "[PASSED] L1StandardBridgeProxy $L1StandardBridgeProxy"

L1ERC721BridgeProxy=`getContractAddress $PathToDeployment/L1ERC721BridgeProxy.json`
echo -n "    Checking L1ERC721BridgeProxy.admin() === L1UpgradeKeyAddress ..."
if [[ xx"$L1UpgradeKeyAddress" == xx`callAddressGetter $L1ERC721BridgeProxy admin` ]]; then
  echo YES
else
  echo NO;
  # exit -1
fi
echo "[PASSED] L1ERC721BridgeProxy $L1ERC721BridgeProxy"

SystemConfigProxy=`getContractAddress $PathToDeployment/SystemConfigProxy.json`
echo "[PASSED] SystemConfigProxy $SystemConfigProxy"

L2OutputOracleProxy=`getContractAddress $PathToDeployment/L2OutputOracleProxy.json`
echo -n "    Checking L2OutputOracleProxy.admin() === L1ProxyAdmin ..."
if [[ xx"$L1ProxyAdmin" == xx`callAddressGetter $L2OutputOracleProxy admin` ]]; then
  echo YES
else
  echo NO;
  exit -1
fi
echo "[PASSED] L2OutputOracleProxy $L2OutputOracleProxy"

L2OutputOracle=`getContractAddress $PathToDeployment/L2OutputOracle.json`
echo "[PASSED] L2OutputOracle $L2OutputOracle"

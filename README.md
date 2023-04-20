# my_project

2022 부산 블록체인 교육 미니프로젝트 

## 선행조건 
- hyperledger fabric 2.2 LTS 설치
- docker, docker-compose 설치
- jq 설치

## 네트워크 수행

1. 네트워크 초기화 
```shell
docker rm -f $(docker ps -aq)
docker rmi -f $(docker images dev-* -q)
docker network prune
docker volume prune
```

2. 네트워크 수행 
```shell
cd 프로젝트 설치 위치/my_project/my_network
./startnetwork.sh
./createchannel.sh
./setAnchorPeerUpdate.sh
./ccp-generate.sh
```

3. 환경설정 및 채널 확인
    해당 내용은 ./createchannel.sh 안에 있음
```shell
export FABRIC_CFG_PATH=${PWD}/config
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer channel list
```

## 체인코드 배포
1. deployCC.sh 내 변수 지정 및 설정 부분 변경 

```shell
./deployCC.sh 
```
2. 체인코드 확인 
```shell
peer lifecycle chaincode --queryinstalled
peer lifecycle chaincode --querycommitted -C howmuchnet
```

## 어플리케이션 수행
1. 필요한 파일 다운 

```shell
cd 프로젝트 경로/application
npm install
node server.js
```
*브라우저 : 
    네트워크 초기화 했을 시 지갑폴더 지우기 

```shell 
 rm -rf wallet
```

2. connection-org1.json 복사하기 
```shell
cp 프로젝트 설치 위치/my_project/my_network/organizations/peerOrganization/org1.example.com/connection-org1.json ./
```

3. localhost:3000 접속 






//package

package main 

//모듈 포함

import (
	"encoding/json"
	"fmt"
	"time"
	"log"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//객체 구조체 선언 (SMartContract)
type SmartContract struct {
	contractapi.Contract
}

//Car, QueryResult 구조체 정의 
type DRecord struct{

	Gardian   string `json:"gadian"` // 보호자 
	RID  string `json:"receiptid"` // 연수증 정보 
	DiagInfo string `json:"diagnosisinfo"` // 진료기록 (진료일, 증상 등)
	PrognoInfo string `json:"prognosisinfo"` //진료 예후 
	Status  string `json:"status"` // registered, verified-verifier (검증여부)
	Price 	int `json:"price"` //진료비 

}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *DRecord
}

//history 결과저장을 위한 구조체
type HistoryQueryResult struct{
	Record *DRecord `json:"record"`
	TxId string `json:"txid"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete bool `json:"isDelete"`
}

//1. Receipt_register : 진료기록 생성 
func (s *SmartContract) Receipt_register(ctx contractapi.TransactionContextInterface, rid string, gardian string, diaginfo string, price int) error {
	
	drecord := DRecord{
		Gardian:   gardian,
		RID:  rid,
		DiagInfo: diaginfo,
		PrognoInfo: "norecord",
		Price :  price,
		Status : "registered",
	}

	dAsBytes, _ := json.Marshal(drecord) // 구조체 정보를 json형식을 ㅗ변환 시키기 마샬링 (직렬화 serialization)

	return ctx.GetStub().PutState(rid, dAsBytes) // 저장 - >인도싱 한다 -> 보증 서명 -> odere -> 원장 동기화
}

//2. Query_record : 생성된 진료 기록 조회 (RID로 영수증 번호로 조회)

func (s *SmartContract) Query_record(ctx contractapi.TransactionContextInterface, rid string) (*DRecord, error) {
	dAsBytes, err := ctx.GetStub().GetState(rid) //읽어오는 것

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if dAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", rid)
	}

	diag := new(DRecord)
	_ = json.Unmarshal(dAsBytes, diag) // json 분서를 구조체로 변환 언마샬링 (객체화)

	return diag, nil
	//히스토리를 읽어오는 것은 포함 완됨 -> marbies_chaincode.go 692line 키가 어떻게 변경되어져 왔ㄴ느지 히스토리를 볼 수 있음
}

//3. Verify_receipt : 영수증에 대한 검증 확인 

func (s *SmartContract) Verify_receipt(ctx contractapi.TransactionContextInterface, rid string, verifier string) error {
	diag, err := s.Query_record(ctx, rid)

	if err != nil {
		return err
	}

	if diag.Status == "registered" {
		diag.Status = verifier

		dAsBytes, _ := json.Marshal(diag)
		return ctx.GetStub().PutState(rid, dAsBytes)
	}else {
		return fmt.Errorf("The receipt is not in ready state")
	}

}

//4. Progno_receipt : 예후등록 

func (s *SmartContract) Progno_receipt(ctx contractapi.TransactionContextInterface, rid string, prognoinfo string) error {
	diag, err := s.Query_record(ctx, rid)

	if err != nil {
		return err
	}

	if diag.PrognoInfo == "norecord" {
		diag.PrognoInfo = prognoinfo
	}else{
		diag.PrognoInfo = diag.PrognoInfo + ", " + prognoinfo
	}

	dAsBytes, _ := json.Marshal(diag)
	return ctx.GetStub().PutState(rid, dAsBytes)

}

//5. GetHistory 기록 가져오기 
func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, rid string) ([]HistoryQueryResult, error) {
	log.Printf("GetHistory: ID %v", rid)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(rid)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var diag DRecord
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &diag)
			if err != nil {
				return nil, err
			}
		} else {
			diag = DRecord {
				RID: rid,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult {
			TxId: response.TxId,
			Timestamp: timestamp,
			Record: &diag,
			IsDelete: response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil 
}

//6. main 
func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create howdog chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting howdog chaincode: %s", err.Error())
	}
}
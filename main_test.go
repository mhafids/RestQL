package main

import (
	"encoding/json"
	queryrestapi "test/QueryRestAPI"
	"testing"
)

type Models struct {
	By        string `json:"by"`
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

func TestCCOperatorJson(t *testing.T) {
	mts := queryrestapi.NewMongoModel()
	var operatorJson = `{"CC":{"phone":{"$not":{"$gt":"25"}}}}`

	var operatorMap map[string]interface{}
	json.Unmarshal([]byte(operatorJson), &operatorMap)

	op, err := mts.Query(operatorMap)
	if err != nil {
		t.Error(err)
	}
	t.Log(op)

	qrapi := queryrestapi.NewQueryJSON(&queryrestapi.QueryConfig{})
	for _, value := range op {
		filtering, err := qrapi.Filter(queryrestapi.ListPayload{Filter: value}, Models{})
		if err != nil {
			t.Error(err)
		}
		t.Log(filtering)
	}
}

// func TestOperatorJson(t *testing.T) {
// 	mts := queryrestapi.NewMongoModel()
// 	var operatorJson = `{"$and":[{"by":"tutorials point"},{"title": "MongoDB Overview"}]}`

// 	var operatorMap map[string]interface{}
// 	json.Unmarshal([]byte(operatorJson), &operatorMap)

// 	op := mts.Query(operatorMap)
// 	t.Log(op)

// 	qrapi := queryrestapi.NewQueryJSON(&queryrestapi.QueryConfig{})
// 	t.Log(qrapi.Filter(queryrestapi.ListPayload{Filter: op}, Models{}))
// }

// func TestNonOperatorJson(t *testing.T) {
// 	mts := queryrestapi.NewMongoModel()
// 	var nonoperatorJson = `{"title": "MongoDB Overview"}`

// 	var nonoperatorMap map[string]interface{}
// 	json.Unmarshal([]byte(nonoperatorJson), &nonoperatorMap)

// 	op := mts.Typeofmap(ParamTypeofMap{data: nonoperatorMap})
// 	t.Log(op)

// 	qrapi := queryrest.NewQueryJSON(&queryrest.QueryConfig{})
// 	t.Log(qrapi.Filter(queryrest.ListPayload{Filter: op}, Models{}))
// }

// func TestComplexOperatorJson(t *testing.T) {
// 	mts := queryrestapi.NewMongoModel()
// 	var complexoperatorJson = `{"$and":[{"$or":[{"first_name":"john"},{"last_name":"john"}]},{"phone":"12345678"}]}`

// 	var complexoperatorMap map[string]interface{}
// 	json.Unmarshal([]byte(complexoperatorJson), &complexoperatorMap)

// 	qrapi := queryrest.NewQueryJSON(&queryrest.QueryConfig{})

// 	op := mts.Typeofmap(ParamTypeofMap{data: complexoperatorMap})
// 	t.Log(op)
// 	t.Log(qrapi.Filter(queryrest.ListPayload{Filter: op}, Models{}))
// }

// func TestOperatorGTJson(t *testing.T) {
// 	mts := queryrestapi.NewMongoModel()
// 	var OperatorGTJson = `{"phone":{"$gt":4,"$lt":6}}`

// 	var OperatorGTMAP map[string]interface{}
// 	json.Unmarshal([]byte(OperatorGTJson), &OperatorGTMAP)

// 	op := mts.Typeofmap(ParamTypeofMap{data: OperatorGTMAP})
// 	t.Log(op)

// 	qrapi := queryrest.NewQueryJSON(&queryrest.QueryConfig{})
// 	t.Log(qrapi.Filter(queryrest.ListPayload{Filter: op}, Models{}))
// }

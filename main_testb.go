package main

// func BenchmarkOperatorJson(b *testing.B) {
// 	mts := queryrestapi.NewMongoModel()
// 	var operatorJson = `{"$and":[{"by":"tutorials point"},{"title": "MongoDB Overview"}]}`

// 	var operatorMap map[string]interface{}
// 	json.Unmarshal([]byte(operatorJson), &operatorMap)

// 	operatorJson = ""
// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			op := mts.MongoTranslate(operatorMap)

// 			qrapi := queryrestapi.NewQueryJSON(&queryrestapi.QueryConfig{})
// 			qrapi.Filter(queryrestapi.ListPayload{Filter: op}, Models{})
// 		}
// 	})
// 	operatorMap = nil
// }

// func BenchmarkNonOperatorJson(b *testing.B) {
// 	mts := &MapToStruct{}
// 	var nonoperatorJson = `{"title": "MongoDB Overview"}`

// 	var nonoperatorMap map[string]interface{}
// 	json.Unmarshal([]byte(nonoperatorJson), &nonoperatorMap)

// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			mts.Typeofmap(ParamTypeofMap{data: nonoperatorMap})
// 		}
// 	})
// }

// func BenchmarkComplexOperatorJson(b *testing.B) {
// 	mts := &MapToStruct{}
// 	var complexoperatorJson = `{"$and":[{"$or":[{"first_name":"john"},{"last_name":"john"}]},{"Phone":"12345678"}]}`

// 	var complexoperatorMap map[string]interface{}
// 	json.Unmarshal([]byte(complexoperatorJson), &complexoperatorMap)

// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			mts.Typeofmap(ParamTypeofMap{data: complexoperatorMap})
// 		}
// 	})

// }

// func BenchmarkOperatorGTJson(b *testing.B) {
// 	mts := &MapToStruct{}
// 	var OperatorGTJson = `{"b":{"$gt":4,"$lt":6}}`

// 	var OperatorGTMAP map[string]interface{}
// 	json.Unmarshal([]byte(OperatorGTJson), &OperatorGTMAP)

// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			mts.Typeofmap(ParamTypeofMap{data: OperatorGTMAP})
// 		}
// 	})
// }

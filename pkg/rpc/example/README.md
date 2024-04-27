Log from the protocol traffic, produced by [main.go](main.go)

```
 client > ["api"]
service < ["api"]
service > ["method","method1","method2","method2B","method2S","methodC","string"]
 client < ["method","method1","method2","method2B","method2S","methodC","string"]
 client > ["string"]
service < ["string"]
service > "testApi"
 client < "testApi"
 client > ["method",true,10,"example"]
service < ["method",true,10,"example"]
service > null
 client < null
 client > ["method1",false]
service < ["method1",false]
service > null
 client < null
 client > ["method1",true]
service < ["method1",true]
service > {"error":"example error"}
 client < {"error":"example error"}
 client > ["method2",null]
service < ["method2",null]
service > {"S":"s","I":1}
 client < {"S":"s","I":1}
 client > ["method2",{"S":"example","I":1000}]
service < ["method2",{"S":"example","I":1000}]
service > {"S":"s","I":1,"arg":{"S":"example","I":1000}}
 client < {"S":"s","I":1,"arg":{"S":"example","I":1000}}
 client > ["method2S"]
service < ["method2S"]
service > "testApi"
 client < "testApi"
 client > ["method2B"]
service < ["method2B"]
service > true
 client < true
 client > ["methodC"]
service < ["methodC"]
service > {"S":"","I":100}
 client < {"S":"","I":100}
service > {"S":"","I":101}
 client < {"S":"","I":101}
service > {"S":"","I":102}
 client < {"S":"","I":102}
service > {"S":"","I":103}
 client < {"S":"","I":103}
service > {"S":"","I":104}
 client < {"S":"","I":104}
```

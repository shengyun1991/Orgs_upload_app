package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strings"
)

type AUTH struct{

}

func (this *AUTH) Init(stub shim.ChaincodeStubInterface) peer.Response{
	return shim.Success(nil)
}

func (this *AUTH) Invoke(stub shim.ChaincodeStubInterface) peer.Response{

	// add

	// check
	// 1. 接收用户的真实姓名和身份证号
	// 2. 做数据的查询
	// 3. 组织数据: -> true:false
	function,parameters := stub.GetFunctionAndParameters()
	if function == "add"{
		return add(stub,parameters[0],parameters[1],parameters[2])
	}else if function == "check"{
		return check(stub,parameters[0],parameters[1])
	}

	return shim.Error("No func match!")

}

func add(stub shim.ChaincodeStubInterface,id string,name string,recode string) peer.Response {
	// value -> name:recode
	err := stub.PutState(id,[]byte(name+":"+recode))
	if err!=nil{
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func check(stub shim.ChaincodeStubInterface,id string,name string) peer.Response{
	data,err := stub.GetState(id)
	if err!=nil{
		return shim.Error(err.Error())
	}

	var result string
	// data 不为空,说明 此 id 所对应的数据之前上传过
	if data !=nil{
		// name:recode
		var str string = string(data[:])
		split := strings.Split(str,":")
		if split[0]==name{
			result = "true"
		}else {
			result ="false"
		}

		//ture:flase

		result = result+":"+split[1]

		return shim.Success([]byte(result))
	}

	return shim.Success([]byte("false:false"))
}

func main(){
	err := shim.Start(new(AUTH))
	if err !=nil{
		fmt.Println("shim.Start error!")
	}
}

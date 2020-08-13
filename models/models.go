package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
)

type AuthModel struct{
	client      apitxn.ChannelClient
	chaincodeID string
}

func (this *AuthModel) ChainCodeQuery(function string, args [][]byte) (response []byte, err error) {
	request := apitxn.QueryRequest{this.chaincodeID, function, args}
	return this.client.Query(request)
}

func (this *AuthModel) ChainCodeUpdate(function string, args [][]byte) (response []byte, err error) {
	request := apitxn.ExecuteTxRequest{ChaincodeID: this.chaincodeID, Fcn: function, Args: args}
	id, err := this.client.ExecuteTx(request)

	return []byte(id.ID), err
}

func (this *AuthModel)Close(){
	this.client.Close()
}

func Initialize(channelID, chainCodeId, userId string,conf string) (*AuthModel, error) {
	config := beego.AppConfig.String(conf)
	fmt.Println("!!!!!!!:",config)
	sdk, err := getSDK(config)
	if err != nil {
		return nil, err
	}

	client, err := sdk.NewChannelClient(channelID, userId)
	if err != nil {
		return nil, err
	}
	return &AuthModel{client, chainCodeId}, nil
}

func getSDK(conf string) (*fabapi.FabricSDK, error) {
	options := fabapi.Options{ConfigFile: conf}
	sdk, err := fabapi.NewSDK(options)
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}

	return sdk, nil
}
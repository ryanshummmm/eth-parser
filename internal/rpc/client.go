package rpc

import (
	"bytes"
	"encoding/json"
	"eth-parser/common"
	"eth-parser/pkg/models"
	"eth-parser/pkg/utils"
	"fmt"
	"io"
	"net/http"
)

func GetLatestBlockNumber() (int, error) {
	response, err := jsonRPCCall(common.EthBlockNumber, []interface{}{})
	if err != nil {
		return 0, err
	}

	blockHex := response.Result.(string)
	return utils.HexToInt(blockHex)
}

func GetBlockByNumber(blockNumber int) (models.Block, error) {
	blockHex := fmt.Sprintf("0x%x", blockNumber)
	response, err := jsonRPCCall(common.EthGetBlockByNumber, []interface{}{blockHex, true})
	if err != nil {
		return models.Block{}, err
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		return models.Block{}, fmt.Errorf("invalid response result format")
	}

	// 将result序列化为JSON
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return models.Block{}, err
	}

	// 将JSON解析为Block结构
	var block models.Block
	err = json.Unmarshal(resultBytes, &block)
	if err != nil {
		return models.Block{}, err
	}

	return block, nil
}

func jsonRPCCall(method string, params []interface{}) (models.JSONRPCResponse, error) {

	request := models.JSONRPCRequest{
		JsonRPC: common.JsonRpcVersion,
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return models.JSONRPCResponse{}, fmt.Errorf("[jsonRPCCall] request body wrong, err=%v", err)
	}

	resp, err := http.Post(common.CloudFlareRpcUrl, common.ApplicationJsonContentType, bytes.NewBuffer(requestBody))
	if err != nil {
		return models.JSONRPCResponse{}, fmt.Errorf("[jsonRPCCall] get response wrong, err=%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.JSONRPCResponse{}, fmt.Errorf("[jsonRPCCall] io wrong, err=%v", err)
	}

	var response models.JSONRPCResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return models.JSONRPCResponse{}, fmt.Errorf("[jsonRPCCall] unmarshal wrong, err=%v", err)
	}

	return response, nil
}

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
	"time"
)

var (
	httpClient *http.Client
)

func init() {
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}

func GetLatestBlockNumber() (int64, error) {
	response, err := jsonRPCCall(common.EthBlockNumber, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block number: %w", err)
	}

	blockHex := response.Result.(string)
	return utils.HexToInt(blockHex)
}

func GetBlockByNumber(blockNumber int64) (models.Block, error) {
	blockHex := fmt.Sprintf("0x%x", blockNumber)
	response, err := jsonRPCCall(common.EthGetBlockByNumber, []interface{}{blockHex, true})
	if err != nil {
		return models.Block{}, fmt.Errorf("failed to get block %d: %w", blockNumber, err)
	}

	resultBytes, err := json.Marshal(response.Result)
	if err != nil {
		return models.Block{}, fmt.Errorf("failed to marshal block %d: %w", blockNumber, err)
	}

	var block models.Block
	if err := json.Unmarshal(resultBytes, &block); err != nil {
		return models.Block{}, fmt.Errorf("failed to unmarshal block %d: %w", blockNumber, err)
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

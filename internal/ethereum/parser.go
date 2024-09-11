package ethereum

import (
	"eth-parser/internal/rpc"
	"eth-parser/internal/storage"
	"eth-parser/pkg/models"
	"fmt"
	"log"
	"sync"
	"time"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int64

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []models.Transaction

	GetSubscribeList() []string
	Unsubscribe(address string) bool
	Start()
	Stop()
}

type EthParser struct {
	storage   storage.Storage
	stopCh    chan struct{}
	logger    *log.Logger
	batchSize int64
}

func NewEthParser(storage storage.Storage, logger *log.Logger) *EthParser {
	return &EthParser{
		storage:   storage,
		stopCh:    make(chan struct{}),
		logger:    logger,
		batchSize: 10, // Process 10 blocks concurrently
	}
}

func (ep *EthParser) GetCurrentBlock() int64 {
	return ep.storage.GetCurrentBlock()
}

func (ep *EthParser) SetCurrentBlock(number int64) {
	ep.storage.SetCurrentBlock(number)
}

func (ep *EthParser) GetSubscribeList() []string {
	return ep.storage.GetSubscribeList()
}

func (ep *EthParser) Subscribe(address string) bool {
	success := ep.storage.Subscribe(address)
	ep.logger.Printf("Subscribed address: %s, success: %v", address, success)
	return success
}

func (ep *EthParser) Unsubscribe(address string) bool {
	success := ep.storage.Unsubscribe(address)
	ep.logger.Printf("Unsubscribed address: %s, success: %v", address, success)
	return success
}

func (ep *EthParser) GetTransactions(address string) []models.Transaction {
	return ep.storage.GetTransactions(address)
}

func (ep *EthParser) updateAndParseBlocks() error {
	latestBlock, err := rpc.GetLatestBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %w", err)
	}

	currentBlock := ep.storage.GetCurrentBlock()
	ep.logger.Printf("Updating blocks from %d to %d", currentBlock, latestBlock)

	for i := currentBlock; i <= latestBlock; i += ep.batchSize {
		end := i + ep.batchSize
		if end > latestBlock {
			end = latestBlock + 1
		}
		if err := ep.processBatch(i, end); err != nil {
			return fmt.Errorf("failed to process batch %d-%d: %w", i, end-1, err)
		}
	}

	for i := currentBlock; i <= latestBlock; i++ {
		block, err := rpc.GetBlockByNumber(i)
		if err != nil {
			return err
		}
		for _, tx := range block.Transactions {
			if ep.storage.IsSubscribed(tx.From) {
				log.Printf("Detect tx for from address: %v", tx.From)
				ep.storage.AddTransaction(tx)
			} else if ep.storage.IsSubscribed(tx.To) {
				log.Printf("Detect tx for to address: %v", tx.To)
				ep.storage.AddTransaction(tx)
			}
		}
	}
	if currentBlock <= latestBlock {
		ep.storage.SetCurrentBlock(latestBlock + 1)
	}

	return nil
}

func (ep *EthParser) processBatch(start, end int64) error {
	var wg sync.WaitGroup
	errCh := make(chan error, end-start)

	for i := start; i < end; i++ {
		wg.Add(1)
		go func(blockNum int64) {
			defer wg.Done()
			if err := ep.processBlock(blockNum); err != nil {
				errCh <- fmt.Errorf("failed to process block %d: %w", blockNum, err)
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return <-errCh
	}
	return nil
}

func (ep *EthParser) processBlock(blockNum int64) error {
	block, err := rpc.GetBlockByNumber(blockNum)
	if err != nil {
		return fmt.Errorf("failed to get block %d: %w", blockNum, err)
	}
	ep.logger.Printf("Processing block %d, transactions: %d", blockNum, len(block.Transactions))

	for _, tx := range block.Transactions {
		if ep.storage.IsSubscribed(tx.From) || ep.storage.IsSubscribed(tx.To) {
			ep.storage.AddTransaction(tx)
			ep.logger.Printf("Detected transaction: from %s to %s, value: %s", tx.From, tx.To, tx.Value)
		}
	}

	return nil
}

func (ep *EthParser) backgroundTask() {
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ep.updateAndParseBlocks(); err != nil {
				ep.logger.Printf("Error updating and parsing blocks: %v", err)
			}
		case <-ep.stopCh:
			ep.logger.Println("Stopping background task")
			return
		}
	}
}

func (ep *EthParser) Start() {
	latestBlockNumber, err := rpc.GetLatestBlockNumber()
	if err != nil {
		ep.logger.Fatalf("Failed to get latest block number: %v", err)
	}
	ep.storage.SetCurrentBlock(latestBlockNumber)
	ep.logger.Printf("Starting parser, initial block: %d", latestBlockNumber)
	go ep.backgroundTask()
}

func (ep *EthParser) Stop() {
	ep.logger.Println("Stopping parser")
	close(ep.stopCh)
}

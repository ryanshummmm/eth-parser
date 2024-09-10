package ethereum

import (
	"eth-parser/internal/rpc"
	"eth-parser/internal/storage"
	"eth-parser/pkg/models"
	"log"
	"time"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

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
	storage storage.Storage
	stopCh  chan struct{}
}

func NewEthParser(storage storage.Storage) *EthParser {
	return &EthParser{
		storage: storage,
		stopCh:  make(chan struct{}),
	}
}

func (ep *EthParser) GetCurrentBlock() int {
	return ep.storage.GetCurrentBlock()
}

func (ep *EthParser) SetCurrentBlock(number int) {
	ep.storage.SetCurrentBlock(number)
}

func (ep *EthParser) GetSubscribeList() []string {
	return ep.storage.GetSubscribeList()
}

func (ep *EthParser) Subscribe(address string) bool {
	return ep.storage.Subscribe(address)
}

func (ep *EthParser) Unsubscribe(address string) bool {
	return ep.storage.Unsubscribe(address)
}

func (ep *EthParser) GetTransactions(address string) []models.Transaction {
	return ep.storage.GetTransactions(address)
}

func (ep *EthParser) updateAndParseBlocks() error {
	latestBlock, err := rpc.GetLatestBlockNumber()
	if err != nil {
		return err
	}

	currentBlock := ep.storage.GetCurrentBlock()
	for i := currentBlock; i <= latestBlock; i++ {
		block, err := rpc.GetBlockByNumber(i)
		if err != nil {
			return err
		}
		log.Printf("block number=%v, tx amount=%v\n", i, len(block.Transactions))

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

func (ep *EthParser) backgroundTask() {
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ep.updateAndParseBlocks(); err != nil {
				log.Printf("Error updating and parsing blocks: %v", err)
			}
		case <-ep.stopCh:
			return
		}
	}
}

func (ep *EthParser) Start() {
	latestBlockNumber, err := rpc.GetLatestBlockNumber()
	if err != nil {
		log.Fatal("end : get block failed, err=", err)
	}
	ep.SetCurrentBlock(latestBlockNumber)
	go ep.backgroundTask()
}

func (ep *EthParser) Stop() {
	close(ep.stopCh)
}

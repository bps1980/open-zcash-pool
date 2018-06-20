package proxy

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
	"sync"

	"github.com/jkkgbe/open-zcash-pool/merkleTree"
	"github.com/jkkgbe/open-zcash-pool/util"
)

const maxBacklog = 3

type heightDiffPair struct {
	diff   *big.Int
	height uint64
}

type Transaction struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	Fee  int    `json:"fee"`
}

type CoinbaseTxn struct {
	Data           string `json:"data"`
	Hash           string `json:"hash"`
	FoundersReward int    `json:"foundersreward"`
}

type BlockTemplate struct {
	sync.RWMutex
	Version       uint32        `json:"version"`
	PrevBlockHash string        `json:"previousblockhash"`
	Transactions  []Transaction `json:"transactions"`
	CoinbaseTxn   CoinbaseTxn   `json:"coinbasetxn"`
	LongpollId    string        `json:"longpollid"`
	Target        string        `json:"target"`
	MinTime       int           `json:"mintime"`
	NonceRange    string        `json:"noncerange"`
	SigOpLimit    int           `json:"sigoplimit"`
	SizeLimit     int           `json:"sizelimit"`
	CurTime       uint32        `json:"curtime"`
	Bits          string        `json:"bits"`
	Height        int           `json:"height"`
}

type Work struct {
	JobId              string
	Version            uint32
	PrevHashReversed   []byte
	MerkleRootReversed []byte
	ReservedField      []byte
	Time               string
	Bits               string
	CleanJobs          bool
	Nonce              string
	SolutionSize       [3]byte
	Solution           [1344]byte
	Header             [4 + 32 + 32 + 32 + 4 + 4 + 32 + 3 + 1344]byte
}

// func (b Work) Difficulty() *big.Int     { return b.difficulty }
// func (b Work) HashNoNonce() common.Hash { return b.hashNoNonce }
// func (b Work) Nonce() uint64            { return b.nonce }
// func (b Work) MixDigest() common.Hash   { return b.mixDigest }
// func (b Work) NumberU64() uint64        { return b.number }

func (s *ProxyServer) fetchWork() {
	rpc := s.rpc()
	t := s.currentWork()
	var reply BlockTemplate
	err := rpc.GetBlockTemplate(&reply)
	if err != nil {
		log.Printf("Error while refreshing block template on %s: %s", rpc.Name, err)
		return
	}
	// No need to update, we have fresh job
	if t != nil && util.BytesToHex(util.ReverseBuffer(t.PrevHashReversed)) == reply.PrevBlockHash {
		return
	}

	// generatedTxHash := CreateRawTransaction(inputs, outputs).TxHash()
	txHashes := make([][32]byte, len(reply.Transactions)+1)
	// txHashes[0] = util.ReverseBuffer(generatedTxHash)
	copy(txHashes[0][:], util.HexToBytes(reply.CoinbaseTxn.Hash)[:32])
	for i, transaction := range reply.Transactions {
		copy(txHashes[i+1][:], util.HexToBytes(transaction.Hash)[:32])
	}

	mtBottomRow := txHashes
	mt := merkleTree.NewMerkleTree(mtBottomRow)
	mtr := mt.MerkleRoot()

	newWork := Work{
		JobId:              "1",
		Version:            reply.Version,
		PrevHashReversed:   util.ReverseBuffer(util.HexToBytes(reply.PrevBlockHash)),
		MerkleRootReversed: util.ReverseBuffer(mtr[:]),
		ReservedField:      util.HexToBytes("0000000000000000000000000000000000000000000000000000000000000000"),
		Time:               util.BytesToHex(util.PackUInt32LE(reply.CurTime)),
		Bits:               util.BytesToHex(util.ReverseBuffer(util.HexToBytes(reply.Bits))),
		CleanJobs:          true,
	}

	// // Copy job backlog and add current one
	// newBlock.headers[reply[0]] = heightDiffPair{
	// 	diff:   util.TargetHexToDiff(reply[2]),
	// 	height: height,
	// }
	// if t != nil {
	// 	for k, v := range t.headers {
	// 		if v.height > height-maxBacklog {
	// 			newBlock.headers[k] = v
	// 		}
	// 	}
	// }
	s.work.Store(&newWork)
	log.Printf("New block to mine on %s at height %d", rpc.Name, reply.Height)

	// Stratum
	if s.config.Proxy.Stratum.Enabled {
		go s.broadcastNewJobs()
	}
}

func (w *Work) BuildHeader(nTime, nBits uint32, noncePart1, noncePart2 []byte) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	_ = binary.Write(buffer, binary.BigEndian, w.Version)
	_, _ = buffer.Write(w.PrevHashReversed)
	_, _ = buffer.Write(w.MerkleRootReversed)
	_, _ = buffer.Write(w.ReservedField)
	_ = binary.Write(buffer, binary.BigEndian, w.Time)
	_ = binary.Write(buffer, binary.BigEndian, w.Bits)
	_, _ = buffer.Write(noncePart1)
	_, _ = buffer.Write(noncePart2)

	return buffer
}

func (w *Work) CreateJob() []interface{} {
	return []interface{}{
		w.JobId,
		util.BytesToHex(util.PackUInt32LE(w.Version)),
		util.BytesToHex(w.PrevHashReversed),
		util.BytesToHex(w.MerkleRootReversed),
		util.BytesToHex(w.ReservedField),
		w.Time,
		w.Bits,
		w.CleanJobs,
	}
}

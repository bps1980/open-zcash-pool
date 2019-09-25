package util

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var Zcash = math.BigPow(10, 8)

var PowLimitMain = new(big.Int).Sub(math.BigPow(2, 243), big.NewInt(1))
var PowLimitTest = new(big.Int).Sub(math.BigPow(2, 251), big.NewInt(1))
var pow256 = math.BigPow(2, 256)

const FirstRewardHalvingBlock = 840000
const pre840Reward int64 = 1000000000
const post840Reward int64 = 625000000

var tAddressPattern = regexp.MustCompile("^t[0-9a-zA-Z]{34}$")
var loginPattern = regexp.MustCompile("^[[:alnum:]]{1,40}$")

var FoundersRewardAddressChangeInterval = 17709.3125
var TestFoundersRewardAddresses = [48]string{
	"t1dBdkjbUMMFQDzDhDgXChN2xeE5XLvyoSW",
	}

func IsValidtAddress(s string) bool {
	return tAddressPattern.MatchString(s)
}

func IsValidLogin(s string) bool {
	return loginPattern.MatchString(s)
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTargetHex(diff int64) string {
	var result [32]uint8
	difficulty := big.NewInt(diff)
	bytes := new(big.Int).Div(PowLimitTest, difficulty).Bytes()
	copy(result[len(result)-len(bytes):], bytes)

	return BytesToHex(result[:])
}

func TargetHexToDiff(targetHex string) *big.Int {
	targetBytes := common.FromHex(targetHex)
	return new(big.Int).Div(pow256, new(big.Int).SetBytes(targetBytes))
}

func ToHex(n int64) string {
	return "0x0" + strconv.FormatInt(n, 16)
}

func FormatReward(reward *big.Int) string {
	return reward.String()
}

func FormatRatReward(reward *big.Rat) string {
	zatoshi := new(big.Rat).SetInt(Zcash)
	reward = reward.Quo(reward, zatoshi)
	return reward.FloatString(8)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}

func String2Big(num string) *big.Int {
	n := new(big.Int)
	n.SetString(num, 0)
	return n
}

func ReverseBuffer(buffer []byte) []byte {
	for i, j := 0, len(buffer)-1; i < j; i, j = i+1, j-1 {
		buffer[i], buffer[j] = buffer[j], buffer[i]
	}
	return buffer
}

func HexToBytes(hexString string) []byte {
	result, _ := hex.DecodeString(hexString)
	return result
}

func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func PackUInt16LE(num uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, num)
	return b
}

func PackUInt32LE(num uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, num)
	return b
}

func PackUInt64LE(num uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return b
}

func PackUInt16BE(num uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, num)
	return b
}

func PackUInt32BE(num uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, num)
	return b
}

func PackUInt64BE(num uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	return b
}

func ReverseUInt32(x uint32) uint32 {
	return (uint32(x)&0xff000000)>>24 |
		(uint32(x)&0x00ff0000)>>8 |
		(uint32(x)&0x0000ff00)<<8 |
		(uint32(x)&0x000000ff)<<24
}

func ReverseHex(hex string) string {
	return BytesToHex(ReverseBuffer(HexToBytes(hex)))
}

func readHex(s string, n int) ([]byte, error) {
	if len(s) > 2*n {
		return nil, errors.New("value oversized")
	}

	bytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	if len(bytes) != n {
		// Pad with zeros
		buf := make([]byte, n)
		copy(buf[n-len(bytes):], bytes)
		buf = bytes
	}

	return bytes, nil
}

func HexToUInt32(s string) uint32 {
	data, err := readHex(s, 4)
	if err != nil {
		return 0
	}

	return binary.BigEndian.Uint32(data)
}

func Sha256d(decrypted []byte) [32]byte {
	round1 := sha256.Sum256(decrypted)
	return sha256.Sum256(round1[:])
}

func GetConstReward(height int64) *big.Int {
	if height >= FirstRewardHalvingBlock {
		return big.NewInt(post840Reward)
	}
	return big.NewInt(pre840Reward)
}
func CreateExtraNonceCounter(seed uint32) uint32 {
	return seed << 27
}

func GetHexTimestamp() string {
	return BytesToHex([]byte(time.Now().String()))
}

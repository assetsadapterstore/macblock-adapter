package macblock

import (
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/common/file"
	"path/filepath"
	"strings"
)

const (
	//币种
	Symbol    = "MAT"
	CurveType = owcrypt.ECC_CURVE_SECP256K1
)


type WalletConfig struct {

	//币种
	Symbol string
	//区块链数据文件
	BlockchainFile string
	//本地数据库文件路径
	dbPath string
	//固定地址
	tokenAddress string
	// 远程服务
	serverAPI string
	//数据目录
	DataDir string
	//本地数据库文件路径
	DBPath string
}

func NewConfig() *WalletConfig {

	c := WalletConfig{}

	//币种
	c.Symbol = Symbol

	//区块链数据文件
	c.BlockchainFile = "blockchain.db"
	//本地数据库文件路径
	c.dbPath = filepath.Join("data", strings.ToLower(c.Symbol), "db")

	//创建目录
	file.MkdirAll(c.dbPath)

	return &c
}

//创建文件夹
func (wc *WalletConfig) makeDataDir() {

	if len(wc.DataDir) == 0 {
		//默认路径当前文件夹./data
		wc.DataDir = "data"
	}

	//本地数据库文件路径
	wc.DBPath = filepath.Join(wc.DataDir, strings.ToLower(wc.Symbol), "db")

	//创建目录
	file.MkdirAll(wc.DBPath)
}
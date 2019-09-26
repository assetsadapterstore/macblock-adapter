package openwtester

import (
	"github.com/assetsadapterstore/macblock-adapter/macblock"
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"path/filepath"
	"testing"
)

var (
	testApp        = "macblock-adapter"
	configFilePath = filepath.Join("conf")
)

var (
	tw *macblock.WalletManager
)

func init() {
	tw = testNewWalletManager()
}

func testNewWalletManager() *macblock.WalletManager {
	wm := macblock.NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "MAT.ini")
	//log.Debug("absFile:", absFile)
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	return wm
}

func TestWalletManager_CreateWallet(t *testing.T) {
	keydir := filepath.Join(tw.Config.DataDir, "key")
	wallet, filePath, err := tw.CreateNewWallet(keydir, "john", "1234qwer")
	if err != nil {
		t.Errorf("CreateNewWallet failed unexpected error: %v\n", err)
		return
	}
	log.Infof("wallet: %+v", wallet)
	log.Infof("keyPath: %s", filePath)

}

func TestWalletManager_GetWalletInfo(t *testing.T) {

	keyFile := filepath.Join(tw.Config.DataDir, "key", "hello-MACcaf763e4780EMgCOUFAHUFCRRgA.key")
	wallet, err := tw.GetWalletInfo(keyFile, "1234qwer")
	if err != nil {
		t.Errorf("GetWalletInfo failed unexpected error: %v\n", err)
		return
	}
	log.Infof("wallet: %+v", wallet)
}

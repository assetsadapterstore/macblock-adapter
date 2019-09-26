/*
 * Copyright 2019 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package macblock

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"path/filepath"
	"testing"
)

var (
	tw *WalletManager
)

func init() {

	tw = testNewWalletManager()
}

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "MAT.ini")
	//log.Debug("absFile:", absFile)
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	//wm.ExplorerClient.Debug = false
	wm.client.Debug = true
	return wm
}

func TestWalletManager_GetAssetBalanceAds(t *testing.T) {
	balance, err := tw.GetAssetBalanceAds("MACja4a7fbe76dBwVUBYFAWZVUWNlA")
	if err != nil {
		t.Errorf("GetAssetBalanceAds failed unexpected error: %v\n", err)
		return
	}
	log.Infof("balance: %s", balance.String())
}

func TestWalletManager_Macpwdencode(t *testing.T) {
	signed := tw.Macpwdencode("1234qwer")
	log.Infof("signed: %s", signed)
}

func TestWalletManager_CreateNewAddress(t *testing.T) {
	address, err := tw.CreateNewAddress("1234qwer")
	if err != nil {
		t.Errorf("CreateNewAddress failed unexpected error: %v\n", err)
		return
	}
	log.Infof("address: %s", address)
}

func TestWalletManager_SignBorn(t *testing.T) {
	signed := tw.SignBorn("", "", "1234qwer")

	log.Infof("signed: %s", signed)
}

func TestWalletManager_CreateNewWallet(t *testing.T) {
	keydir := filepath.Join(tw.Config.DataDir, "key")
	wallet, filePath, err := tw.CreateNewWallet(keydir, "kelly", "1234qwer")
	if err != nil {
		t.Errorf("CreateNewWallet failed unexpected error: %v\n", err)
		return
	}
	log.Infof("wallet: %+v", wallet)
	log.Infof("keyPath: %s", filePath)
}

func TestWalletManager_GetWalletInfo(t *testing.T) {
	keyFile := filepath.Join(tw.Config.DataDir, "key", "john-MACx6150b0728bVdQDOAABCYFAUN1U.key")
	wallet, err := tw.GetWalletInfo(keyFile, "1234qwer")
	if err != nil {
		t.Errorf("GetWalletInfo failed unexpected error: %v\n", err)
		return
	}
	log.Infof("wallet: %+v", wallet)
}

func TestWalletManager_SendTransaction(t *testing.T) {

	rawTx := &openwallet.RawTransaction{
		Coin: openwallet.Coin{
			Symbol:     "MAT",
			IsContract: false,
		},
		To: map[string]string{"MACja4a7fbe76dBwVUBYFAWZVUWNlA": "0.01"},
	}

	rawTx.SetExtParam("memo", "john")

	keyFile := filepath.Join(tw.Config.DataDir, "key", "john-MACx6150b0728bVdQDOAABCYFAUN1U.key")

	wallet, err := tw.GetWalletInfo(keyFile, "1234qwer")
	if err != nil {
		t.Errorf("SendTransaction failed unexpected error: %v\n", err)
		return
	}

	tx, err := tw.SendTransaction(wallet, "1234qwer", rawTx)
	if err != nil {
		t.Errorf("SendTransaction failed unexpected error: %v\n", err)
		return
	}
	log.Infof("tx: %+v", tx)
}

func TestWalletManager_GetBlockHeight(t *testing.T) {
	height, err := tw.GetBlockHeight()
	if err != nil {
		t.Errorf("GetBlockHeight failed unexpected error: %v\n", err)
		return
	}
	log.Infof("height: %d", height)
}

func TestWalletManager_GetTransactionRecordHight(t *testing.T) {
	block, err := tw.GetTransactionRecordHight(338567)
	if err != nil {
		t.Errorf("GetTransactionRecordHight failed unexpected error: %v\n", err)
		return
	}
	log.Infof("block: %+v", block)
	for _, tx := range block.txDetails {
		log.Infof("tx: %+v", tx)
	}
}
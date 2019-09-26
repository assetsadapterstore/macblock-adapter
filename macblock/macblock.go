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
	"github.com/asdine/storm"
	"github.com/astaxie/beego/config"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"path/filepath"
)

//FullName 币种全名
func (wm *WalletManager) FullName() string {
	return "cxcblock"
}

//CurveType 曲线类型
func (wm *WalletManager) CurveType() uint32 {
	return owcrypt.ECC_CURVE_SECP256K1
}

//Symbol 币种标识
func (wm *WalletManager) Symbol() string {
	return Symbol
}

//小数位精度
func (wm *WalletManager) Decimal() int32 {
	return 8
}

//AddressDecode 地址解析器
func (wm *WalletManager) GetAddressDecode() openwallet.AddressDecoder {
	return nil
}

//TransactionDecoder 交易单解析器
func (wm *WalletManager) GetTransactionDecoder() openwallet.TransactionDecoder {
	return nil
}

//GetBlockScanner 获取区块链
func (wm *WalletManager) GetBlockScanner() openwallet.BlockScanner {
	return wm.Blockscanner
}

//LoadAssetsConfig 加载外部配置
func (wm *WalletManager) LoadAssetsConfig(c config.Configer) error {

	wm.Config.serverAPI = c.String("serverAPI")
	wm.Config.tokenAddress = c.String("tokenAddress")
	wm.Config.DataDir = c.String("dataDir")

	wm.client = NewClient(wm.Config.serverAPI, false)

	//数据文件夹
	wm.Config.makeDataDir()

	blockchaindb, err := storm.Open(filepath.Join(wm.Config.dbPath, wm.Config.BlockchainFile))
	if err != nil {
		return err
	}

	wm.blockChainDB = blockchaindb

	return nil
}

//InitAssetsConfig 初始化默认配置
func (wm *WalletManager) InitAssetsConfig() (config.Configer, error) {
	return nil, nil
}

//GetAssetsLogger 获取资产账户日志工具
func (wm *WalletManager) GetAssetsLogger() *log.OWLogger {
	return wm.Log
}

//GetSmartContractDecoder 获取智能合约解析器
func (wm *WalletManager) GetSmartContractDecoder() openwallet.SmartContractDecoder {
	return nil
}

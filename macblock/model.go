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
	"fmt"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/tidwall/gjson"
)

type MACWallet struct {
	Alias         string `json:"alias"`
	Address       string `json:"NewTokenAddress" `
	WalletKey     string `json:"WalletKey"`
	MnemonicWords string `json:"MnemonicWords"`
}

// 加密后的MACWallet的JSON结构
type encryptedWalletJSON struct {
	Alias   string     `json:"alias"`
	Address string     `json:"address"`
	Crypto  cryptoJSON `json:"crypto"`
}

// 加密内容的JSON结构
type cryptoJSON struct {
	CipherKey   string `json:"cipherKey"`
	CipherWords string `json:"cipherWords"`
}

type Block struct {
	Hash              string
	Previousblockhash string
	Height            uint64 `storm:"id"`
	Time              uint64
	Fork              bool
	txDetails         []*Transaction
}

func NewBlock(height uint64, json *gjson.Result) *Block {
	obj := &Block{}
	//解析json
	obj.Height = height
	obj.Hash = gjson.Get(json.Raw, "blockhash").String()
	obj.Previousblockhash = gjson.Get(json.Raw, "parenthash").String()
	obj.Time = gjson.Get(json.Raw, "time").Uint()
	txs := gjson.Get(json.Raw, "Content")
	txDetails := make([]*Transaction, 0)

	if txs.IsArray() {
		for _, tx := range txs.Array() {
			txObj := NewTransaction(&tx)
			txObj.BlockHeight = height
			txObj.BlockHash = obj.Hash
			txDetails = append(txDetails, txObj)
		}
	}
	obj.txDetails = txDetails

	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader(symbol string) *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	obj.Previousblockhash = b.Previousblockhash
	obj.Height = b.Height
	obj.Time = b.Time
	obj.Symbol = symbol

	return &obj
}

type Transaction struct {
	TxID        string
	FromToken   string
	ToToken     string
	Amount      string
	Time        int64
	Note        string
	BlockHeight uint64
	BlockHash   string
}

func NewTransaction(json *gjson.Result) *Transaction {

	obj := Transaction{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "hash").String()
	obj.FromToken = gjson.Get(json.Raw, "fromtoken").String()
	obj.ToToken = gjson.Get(json.Raw, "totoken").String()
	obj.Amount = gjson.Get(json.Raw, "amount").String()
	obj.Time = gjson.Get(json.Raw, "time").Int()
	obj.Note = gjson.Get(json.Raw, "note").String()

	return &obj
}


//UnscanRecords 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	TxID        string
	Reason      string
}

func NewUnscanRecord(height uint64, txID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.TxID = txID
	obj.Reason = reason
	obj.ID = common.Bytes2Hex(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, txID))))
	return &obj
}
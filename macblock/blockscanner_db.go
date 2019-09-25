/*
 * Copyright 2018 The openwallet Authors
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
	"errors"
)

const (
	blockchainBucket  = "blockchain" // blockchain dataset
)

//SaveLocalBlockHead 记录区块高度和hash到本地
func (bs *MACBlockScanner) SaveLocalNewBlock(blockHeight uint64, blockHash string) error {

	//获取本地区块高度

	bs.wm.blockChainDB.Set(blockchainBucket, "blockHeight", &blockHeight)
	bs.wm.blockChainDB.Set(blockchainBucket, "blockHash", &blockHash)

	return nil
}

//GetLocalBlockHead 获取本地记录的区块高度和hash
func (bs *MACBlockScanner) GetLocalNewBlock() (uint64, string) {

	var (
		blockHeight uint64
		blockHash   string
	)

	////获取本地区块高度

	bs.wm.blockChainDB.Get(blockchainBucket, "blockHeight", &blockHeight)
	bs.wm.blockChainDB.Get(blockchainBucket, "blockHash", &blockHash)

	return blockHeight, blockHash
}

//SaveLocalBlock 记录本地新区块
func (bs *MACBlockScanner) SaveLocalBlock(blockHeader *Block) error {


	bs.wm.blockChainDB.Save(blockHeader)

	return nil
}

//GetLocalBlock 获取本地区块数据
func (bs *MACBlockScanner) GetLocalBlock(height uint64) (*Block, error) {

	var (
		blockHeader Block
	)

	err := bs.wm.blockChainDB.One("BlockNumber", height, &blockHeader)
	if err != nil {
		return nil, err
	}

	return &blockHeader, nil
}


//获取未扫记录
func (bs *MACBlockScanner) GetUnscanRecords() ([]*UnscanRecord, error) {

	var list []*UnscanRecord
	err := bs.wm.blockChainDB.All(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

//SaveUnscanRecord 保存交易记录到钱包数据库
func (bs *MACBlockScanner) SaveUnscanRecord(record *UnscanRecord) error {

	if record == nil {
		return errors.New("the unscan record to save is nil")
	}

	////获取本地区块高度

	return bs.wm.blockChainDB.Save(record)
}

//DeleteUnscanRecord 删除指定高度的未扫记录
func (bs *MACBlockScanner) DeleteUnscanRecord(height uint64) error {
	//获取本地区块高度


	var list []*UnscanRecord
	err := bs.wm.blockChainDB.Find("BlockHeight", height, &list)
	if err != nil {
		return err
	}

	for _, r := range list {
		bs.wm.blockChainDB.DeleteStruct(r)
	}

	return nil
}
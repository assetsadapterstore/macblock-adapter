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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/asdine/storm"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/hdkeystore"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/imroc/req"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type WalletManager struct {
	openwallet.AssetsAdapterBase

	Config          *WalletConfig                   // 节点配置
	Decoder         openwallet.AddressDecoder       //地址编码器
	TxDecoder       openwallet.TransactionDecoder   //交易单编码器
	Log             *log.OWLogger                   //日志工具
	ContractDecoder openwallet.SmartContractDecoder //智能合约解析器
	Blockscanner    openwallet.BlockScanner         //区块扫描器
	client          *Client                         //远程客户端
	blockChainDB    *storm.DB                       //区块链数据库
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig()
	wm.Blockscanner = NewMACBlockScanner(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	return &wm
}

// GetAssetBalanceAds 获取余额
func (wm *WalletManager) GetAssetBalanceAds(address string) (decimal.Decimal, error) {

	param := req.Param{
		"action":       "GetAssetBalanceAds",
		"tokenaddress": address,
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return decimal.Zero, nil
	}

	balance := result.Get("AssetBalance").String()
	return decimal.NewFromString(balance)
}

// CreateNewAddress 创建地址
func (wm *WalletManager) CreateNewAddress(password string) (string, error) {

	pwdencrypt := wm.Macpwdencode(password)

	param := req.Param{
		"action":     "IncreaseTokenAddress2",
		"pwdencrypt": pwdencrypt,
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return "", err
	}

	address := result.Get("NewTokenAddress").String()
	return address, nil
}

func (wm *WalletManager) GetmyWalletKey2(address, password string) (string, error) {
	sign := wm.SignBorn("", "", password)

	param := req.Param{
		"action": "GetmyWalletKey2",
		"token":  address,
		"sign":   sign,
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return "", err
	}

	walletKey := result.Get("WalletKey").String()
	return walletKey, nil
}

func (wm *WalletManager) GetMnemonicWords2(address, walletKey, password string) (string, error) {
	sign := wm.SignBorn(walletKey, "", password)

	param := req.Param{
		"action": "GetMnemonicWords2",
		"token":  address,
		"sign":   sign,
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return "", err
	}

	mnemonicWords := result.Get("MnemonicWords").String()
	return mnemonicWords, nil
}

// CreateNewWallet 创建新钱包
func (wm *WalletManager) CreateNewWallet(keydir, alias, password string) (*MACWallet, string, error) {
	address, err := wm.CreateNewAddress(password)
	if err != nil {
		return nil, "", err
	}
	walletKey, err := wm.GetmyWalletKey2(address, password)
	if err != nil {
		return nil, "", err
	}
	mnemonicWords, err := wm.GetMnemonicWords2(address, walletKey, password)
	if err != nil {
		return nil, "", err
	}
	mtsign, err := wm.GetMtsign2(address, walletKey, mnemonicWords, password)
	if err != nil {
		return nil, "", err
	}
	wallet := &MACWallet{
		Alias:         alias,
		Address:       address,
		WalletKey:     walletKey,
		MnemonicWords: mnemonicWords,
		MtSign:        mtsign,
	}

	//加密保存到文件夹
	encryptJSON, err := wm.EncryptWallet(wallet, password)
	if err != nil {
		return nil, "", err
	}

	filePath := filepath.Join(keydir, hdkeystore.KeyFileName(alias, address)+".key")
	err = writeKeyFile(filePath, encryptJSON)
	if err != nil {
		return nil, "", err
	}

	return wallet, filePath, nil
}

func (wm *WalletManager) EncryptWallet(wallet *MACWallet, password string) ([]byte, error) {

	passwordHash := crypto.SHA256([]byte(password))

	cipherKey, err := crypto.AESEncrypt([]byte(wallet.WalletKey), passwordHash)
	if err != nil {
		return nil, err
	}

	cipherWords, err := crypto.AESEncrypt([]byte(wallet.MnemonicWords), passwordHash)
	if err != nil {
		return nil, err
	}

	cipherMtSign, err := crypto.AESEncrypt([]byte(wallet.MtSign), passwordHash)
	if err != nil {
		return nil, err
	}

	cryptoStruct := cryptoJSON{
		CipherKey:    hex.EncodeToString(cipherKey),
		CipherWords:  hex.EncodeToString(cipherWords),
		CipherMtSign: hex.EncodeToString(cipherMtSign),
	}

	encryptedWallet := encryptedWalletJSON{
		Alias:   wallet.Alias,
		Address: wallet.Address,
		Crypto:  cryptoStruct,
	}
	return json.MarshalIndent(encryptedWallet, "", "\t")
}

// GetWalletInfo 通过密钥文件解析钱包
func (wm *WalletManager) GetWalletInfo(keyFile, password string) (*MACWallet, error) {

	keyjson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	wallet, err := wm.DecryptWallet(keyjson, password)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (wm *WalletManager) DecryptWallet(walletJson []byte, password string) (*MACWallet, error) {

	k := new(encryptedWalletJSON)
	if err := json.Unmarshal(walletJson, k); err != nil {
		return nil, err
	}

	passwordHash := crypto.SHA256([]byte(password))
	ck, err := hex.DecodeString(k.Crypto.CipherKey)
	if err != nil {
		return nil, err
	}

	cw, err := hex.DecodeString(k.Crypto.CipherWords)
	if err != nil {
		return nil, err
	}

	cm, err := hex.DecodeString(k.Crypto.CipherMtSign)
	if err != nil {
		return nil, err
	}

	walletKey, err := crypto.AESDecrypt(ck, passwordHash)
	if err != nil {
		return nil, err
	}

	mnemonicWords, err := crypto.AESDecrypt(cw, passwordHash)
	if err != nil {
		return nil, err
	}

	mtSign, err := crypto.AESDecrypt(cm, passwordHash)
	if err != nil {
		return nil, err
	}

	wallet := &MACWallet{
		Alias:         k.Alias,
		Address:       k.Address,
		WalletKey:     string(walletKey),
		MnemonicWords: string(mnemonicWords),
		MtSign:        string(mtSign),
	}

	return wallet, nil
}

func (wm *WalletManager) GetMtsign2(token, walletKey, mnemonicWords, password string) (string, error) {

	sign := wm.SignBorn(walletKey, mnemonicWords, password)

	param := req.Param{
		"action": "GetMtsign2",
		"token":  token,
		"sign":   sign,
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return "", err
	}

	mtsign := result.Get("Mtsign").String()
	return mtsign, nil
}

func (wm *WalletManager) AssetTransferMN2(fromtoken, totoken, amount, note, mtsign, password string) (string, error) {

	sign := wm.SignBorn("", mtsign, password)

	param := req.Param{
		"action":    "AssetTransferMN2",
		"fromtoken": fromtoken,
		"totoken":   totoken,
		"amount":    amount,
		"sign":      sign,
		"note":      note,
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return "", err
	}

	txid := result.Get("TranHash").String()
	return txid, nil
}

func (wm *WalletManager) SendTransaction(wallet *MACWallet, password string, rawTx *openwallet.RawTransaction) (*openwallet.Transaction, error) {

	var (
		fromtoken string
		totoken   string
		toamount  string
		note      string
	)

	fromtoken = wallet.Address

	for to, amount := range rawTx.To {
		totoken = to
		toamount = amount
	}

	note = rawTx.GetExtParam().Get("memo").String()

	balance, err := wm.GetAssetBalanceAds(fromtoken)
	if err != nil {
		return nil, err
	}

	totalAmount, _ := decimal.NewFromString(toamount)

	if balance.LessThan(totalAmount) {
		return nil, openwallet.Errorf(openwallet.ErrInsufficientBalanceOfAddress, "address's balance is not enough")
	}

	txid, err := wm.AssetTransferMN2(fromtoken, totoken, toamount, note, wallet.MtSign, password)
	if err != nil {
		return nil, err
	}

	rawTx.TxID = txid
	rawTx.IsSubmit = true

	txFrom := []string{fmt.Sprintf("%s:%s", fromtoken, toamount)}
	txTo := []string{fmt.Sprintf("%s:%s", totoken, toamount)}
	decimals := wm.Decimal()

	//记录一个交易单
	tx := &openwallet.Transaction{
		From:       txFrom,
		To:         txTo,
		Amount:     rawTx.TxAmount,
		Coin:       rawTx.Coin,
		TxID:       rawTx.TxID,
		Decimal:    decimals,
		Fees:       "0",
		SubmitTime: time.Now().Unix(),
		ExtParam:   rawTx.ExtParam,
	}

	tx.WxID = openwallet.GenTransactionWxID(tx)

	return tx, nil
}

func (wm *WalletManager) Macpwdencode(password string) string {
	a := randSeq(8)
	b := crypto.GetMD5(crypto.GetMD5(password))
	c := b + a
	d := c + b
	e := crypto.GetMD5(d)
	f := e + b + a
	return f
}

func (wm *WalletManager) SignBorn(a, b, c string) string {
	d := crypto.GetMD5(crypto.GetMD5(c))
	e := common.NewString(time.Now().UnixNano() / 1e6).String()
	f := strings.ReplaceAll(a+b+d+e, " ", "")
	f = strings.ToLower(f)
	hash := crypto.SHA256([]byte(f))
	g := hex.EncodeToString(hash)
	return g + e
}

// randSeq 随机字符串
func randSeq(n int) string {

	letters := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//writeKeyFile
func writeKeyFile(file string, content []byte) error {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()
	return os.Rename(f.Name(), file)
}

func (wm *WalletManager) GetBlockHeight() (uint64, error) {

	param := req.Param{
		"action": "GetBlockHeight",
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return 0, err
	}

	height := result.Get("BlockHeight").Uint()
	return height, nil
}

func (wm *WalletManager) GetTransactionRecordHight(height uint64) (*Block, error) {

	param := req.Param{
		"action": "GetTransactionRecordHight",
		"height": height,
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return nil, err
	}

	block := NewBlock(height, result)
	return block, nil
}

func (wm *WalletManager) GetTransactionRecordHash(hash string) (*Transaction, error) {

	param := req.Param{
		"action": "GetTransactionRecordHash",
		"hash":   hash,
	}

	result, err := wm.client.Call(param)
	if err != nil {
		return nil, err
	}

	content := result.Get("Content")

	if content.IsArray() {
		for _, tx := range content.Array() {
			trx := NewTransaction(&tx)
			return trx, nil
		}
	}

	return nil, fmt.Errorf("can not find tx: %s", hash)
}

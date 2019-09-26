# macblock-adapter

macblock-adapter适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 如何集成

1. 先下载其官方手机钱包，并完成注册。
https://www.macblock.io/download.php

2. 登录开发者平台，在【管理首页】把【合约请求地址】记下来
https://www.macblock.io/Developer_center/

3. 新建MAT.ini文件，编辑如下内容：

```ini

# node api url
serverAPI = "http://"

# Cache data file directory, default = "", current directory: ./data
dataDir = ""

```

把【合约地址】填充到serverAPI，请使用https。

4. 集成代码示例

```go

    //创建macblock钱包管理实例
    wm := macblock.NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "MAT.ini")
	//log.Debug("absFile:", absFile)
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	
	//创建钱包
	keydir := filepath.Join(tw.Config.DataDir, "key")
	wallet, filePath, err := tw.CreateNewWallet(keydir, "john", "1234qwer")

    //通过密钥和密码加载钱包	
	wallet, err := tw.GetWalletInfo(keyFile, "1234qwer")

	//指定钱包发起交易
	tx, err := tw.SendTransaction(wallet, "1234qwer", rawTx)

    //获取扫描器	
    scanner := tw.GetBlockScanner()
    //设置查找地址算法
    scanner.SetBlockScanTargetFunc(scanTargetFunc)
    //注册订阅者
    sub := subscriberSingle{}
    scanner.AddObserver(&sub)
    //运行扫描器
    scanner.Run()
    	
```
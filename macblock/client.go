package macblock

import (
	"errors"
	"fmt"
	"github.com/blocktree/openwallet/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

type Client struct {
	BaseURL string
	Debug   bool
	Client  *req.Req
}

func NewClient(url string, debug bool) *Client {
	c := Client{
		BaseURL: url,
		Debug:   debug,
	}

	api := req.New()
	c.Client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *Client) Call(param req.Param) (*gjson.Result, error) {

	if c.Client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	if c.Debug {
		log.Std.Info("Start Request API...")
	}

	r, err := c.Client.Post(c.BaseURL, param)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())
	err = isError(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

//isError 是否报错
func isError(result *gjson.Result) error {
	var (
		err error
	)

	/*
		{
		    "errCode": 1,
		    "Msg": "\u5730\u5740\u6709\u8bef",
		    "AllAsset": "",
		    "AssetBalance": "",
		    "LockedBalance": "",
		    "Assetpaifa": null
		}
	*/

	if result.Get("errCode").Int() == 0 {
		return nil
	}

	errInfo := fmt.Sprintf("[%d]%s",
		result.Get("errCode").Int(),
		result.Get("Msg").String())
	err = errors.New(errInfo)

	return err
}



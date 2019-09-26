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
	"github.com/blocktree/openwallet/log"
	"testing"
)

func TestMACBlockScanner_GetBalanceByAddress(t *testing.T) {

	addrs := []string{
		"MACcaf763e4780EMgCOUFAHUFCRRgA",
	}

	balances, err := tw.GetBlockScanner().GetBalanceByAddress(addrs...)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	for _, b := range balances {
		log.Infof("balance[%s] = %s", b.Address, b.Balance)
		log.Infof("UnconfirmBalance[%s] = %s", b.Address, b.UnconfirmBalance)
		log.Infof("ConfirmBalance[%s] = %s", b.Address, b.ConfirmBalance)
	}
}

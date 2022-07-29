package utils

import "net/http"

var GetNetworkError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot get networks"}}

// var EthError = Response{http.StatusBadGateway, map[string]any{"message": "cannot connect to network"}}

var UnsupportNetworkError = Response{http.StatusNotFound, map[string]any{"message": "network not found"}}

var AddWalletError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot add wallet"}}

var DeleteWalletError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot delete wallet"}}

var FindWalletError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot find wallet"}}

var GetBalanceError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot get balance"}}

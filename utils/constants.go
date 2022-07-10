package utils

import "net/http"

var GetNetworkError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot connect get networks"}}

var EthError = Response{http.StatusBadGateway, map[string]any{"message": "cannot connect to network"}}

var UnsupportNetworkError = Response{http.StatusNotFound, map[string]any{"message": "network not found"}}

var AddWalletError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot add wallet"}}

var DeleteAddressesError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot delete addresses"}}

var FindWalletError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot find wallet"}}

var GetBalanceError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot get balance"}}

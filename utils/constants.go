package utils

import "net/http"

var SecretKey = "Richard Yang"

var GetNetworkError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot connect get networks"}}

var EthError = Response{http.StatusBadGateway, map[string]any{"message": "cannot connect to network"}}

var UnsupportNetworkError = Response{http.StatusNotFound, map[string]any{"message": "network not found"}}

var LoginAuthError = Response{http.StatusBadRequest, map[string]any{"message": "incorrect email or password"}}

var VerificationCodeError = Response{http.StatusBadRequest, map[string]any{"message": "wrong verification code"}}

var UserLoginError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot login"}}

var AddWalletError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot add wallet address"}}

var FindUserError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot find user"}}

var CreateUserError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot create user"}}

var SendCodeError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot send code"}}

var DeleteAddressesError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot delete addresses"}}

var FindWalletError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot find wallet"}}

var ChangePasswordError = Response{http.StatusInternalServerError, map[string]any{"message": "cannot change password"}}

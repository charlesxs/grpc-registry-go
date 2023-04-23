package gclient

import "errors"

var ErrClientInit = errors.New("grpcClient init error")
var ErrConfigNotFound = errors.New("config not found error")
var ErrConfig = errors.New("config error")
var ErrCreateConn = errors.New("create clientConn error")

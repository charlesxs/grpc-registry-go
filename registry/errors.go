package registry

import "errors"

var ErrRegistry = errors.New("registry error")
var ErrRegistryOption = errors.New("registry option error")
var ErrConfigNotFound = errors.New("config not found")

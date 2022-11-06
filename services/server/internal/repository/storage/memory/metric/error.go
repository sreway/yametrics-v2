package metric

import "errors"

var (
	ErrOpenFile     = errors.New("can't open store file")
	ErrStoreMetrics = errors.New("can't store metrics")
	ErrLoadMetrics  = errors.New("can't load metrics")
)

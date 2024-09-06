package storage

import (
	"context"
	sPb "github.com/c12s/scheme/stellar"
)

type DB interface {
	List(context.Context, *sPb.ListReq) (*sPb.ListResp, error)
	Get(context.Context, *sPb.GetReq) (*sPb.GetResp, error)
	StartCollector(context.Context)
}

package syncer

import (
	sPb "github.com/c12s/scheme/stellar"
)

type Syncer interface {
	Sub(f func(msg *sPb.LogBatch))
}

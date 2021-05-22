package app

import "context"

type ContextProvider interface {
	CtxWithTx() (context.Context, context.CancelFunc)
	TxCommit(ctx context.Context) error
}

package context

import (
	"context"
	"fmt"

	"mentat/internal/core"
)

var mentatContextKey struct{}

func SetMentatContext(ctx context.Context, m *core.Mentat) context.Context {
	return context.WithValue(ctx, mentatContextKey, m)
}

func GetMentatContext(ctx context.Context) (*core.Mentat, error) {
	mentat := ctx.Value(mentatContextKey).(*core.Mentat)
	var err error
	if mentat == nil {
		err = fmt.Errorf("could not load mentat instance from context")
	}
	return mentat, err
}

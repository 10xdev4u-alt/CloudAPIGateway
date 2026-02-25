package wasm

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type Engine struct {
	runtime wazero.Runtime
}

func NewEngine(ctx context.Context) *Engine {
	return &Engine{
		runtime: wazero.NewRuntime(ctx),
	}
}

func (e *Engine) Close(ctx context.Context) error {
	return e.runtime.Close(ctx)
}

func (e *Engine) CompileModule(ctx context.Context, binary []byte) (wazero.CompiledModule, error) {
	return e.runtime.CompileModule(ctx, binary)
}

func (e *Engine) InstantiateModule(ctx context.Context, cm wazero.CompiledModule, config wazero.ModuleConfig) (api.Module, error) {
	return e.runtime.InstantiateModule(ctx, cm, config)
}

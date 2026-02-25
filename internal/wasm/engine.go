package wasm

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Engine struct {
	runtime wazero.Runtime
}

func NewEngine(ctx context.Context) *Engine {
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	return &Engine{
		runtime: r,
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

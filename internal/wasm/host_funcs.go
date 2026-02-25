package wasm

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func (m *Manager) registerHostFunctions(ctx context.Context) error {
	_, err := m.engine.runtime.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(m.logInfo).
		Export("log_info").
		Instantiate(ctx)
	return err
}

func (m *Manager) logInfo(ctx context.Context, mod api.Module, ptr uint32, length uint32) {
	if bytes, ok := mod.Memory().Read(ptr, length); ok {
		// Just a simple print for now to establish the ABI
		println("Wasm Plugin Log:", string(bytes))
	}
}

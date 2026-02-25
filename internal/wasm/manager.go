package wasm

import (
	"context"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
)

type Manager struct {
	engine   *Engine
	plugins map[string]wazero.CompiledModule
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		engine:  NewEngine(ctx),
		plugins: make(map[string]wazero.CompiledModule),
	}
}

func (m *Manager) LoadPlugin(ctx context.Context, name, path string) error {
	binary, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read plugin binary: %w", err)
	}

	cm, err := m.engine.CompileModule(ctx, binary)
	if err != nil {
		return fmt.Errorf("failed to compile plugin: %w", err)
	}

	m.plugins[name] = cm
	return nil
}

func (m *Manager) GetPlugin(name string) (wazero.CompiledModule, bool) {
	cm, ok := m.plugins[name]
	return cm, ok
}

func (m *Manager) Close(ctx context.Context) error {
	return m.engine.Close(ctx)
}

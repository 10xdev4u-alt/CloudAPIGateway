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

func (m *Manager) RunGreet(ctx context.Context, name string) error {
	cm, ok := m.plugins[name]
	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Instantiate the module
	mod, err := m.engine.InstantiateModule(ctx, cm, wazero.NewModuleConfig().
		WithName(name).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr))
	if err != nil {
		return fmt.Errorf("failed to instantiate module: %w", err)
	}
	defer mod.Close(ctx)

	greet := mod.ExportedFunction("greet")
	if greet == nil {
		return fmt.Errorf("greet function not found in plugin %s", name)
	}

	_, err = greet.Call(ctx)
	if err != nil {
		return fmt.Errorf("failed to call greet: %w", err)
	}

	return nil
}

func (m *Manager) Close(ctx context.Context) error {
	return m.engine.Close(ctx)
}

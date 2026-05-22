package sigterm

import (
	"fmt"
	"context"
	"sync"
)

// topoSort returns hooks in reverse-dependency order (leaves first).
func topoSort(hooks map[string]*Hook) ([][]string, error) {
	visited := make(map[string]int) // 0=unvisited, 1=in-progress, 2=done
	var order []string

	var visit func(name string) error
	visit = func(name string) error {
		switch visited[name] {
		case 2:
			return nil
		case 1:
			return fmt.Errorf("cyclic dependency detected at %q", name)
		}
		visited[name] = 1
		hook, ok := hooks[name]
		if !ok {
			return fmt.Errorf("unknown hook dependency %q", name)
		}
		for _, dep := range hook.deps {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visited[name] = 2
		order = append(order, name)
		return nil
	}

	for name := range hooks {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	// Group into parallel stages
	position := make(map[string]int, len(order))
	for i, name := range order {
		position[name] = i
	}

	stages := make([][]string, len(order))
	for _, name := range order {
		hook := hooks[name]
		stage := 0
		for _, dep := range hook.deps {
			if s := position[dep] + 1; s > stage {
				stage = s
			}
		}
		stages[stage] = append(stages[stage], name)
	}

	var result [][]string
	for _, s := range stages {
		if len(s) > 0 {
			result = append(result, s)
		}
	}
	return result, nil
}

// runInOrder executes each stage concurrently, waiting for each stage before proceeding.
func runInOrder(ctx context.Context, stages [][]string, hooks map[string]*Hook) error {
	for _, stage := range stages {
		var (
			wg   sync.WaitGroup
			mu   sync.Mutex
			errs []error
		)
		for _, name := range stage {
			wg.Add(1)
			go func(h *Hook) {
				defer wg.Done()
				if err := h.handler(ctx); err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("hook %q: %w", h.name, err))
					mu.Unlock()
				}
			}(hooks[name])
		}
		wg.Wait()
		if len(errs) > 0 {
			return errs[0]
		}
	}
	return nil
}

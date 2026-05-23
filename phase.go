package sigterm

// Phase represents a named shutdown phase grouping handlers
// that should run together before advancing to the next phase.
type Phase struct {
	Name     string
	Priority int
}

// Predefined shutdown phases with conventional priority levels.
var (
	PhasePreShutdown  = Phase{Name: "pre-shutdown", Priority: 100}
	PhaseApplication  = Phase{Name: "application", Priority: 75}
	PhaseInfrastructure = Phase{Name: "infrastructure", Priority: 50}
	PhaseCleanup      = Phase{Name: "cleanup", Priority: 25}
)

// phaseRegistry maps handler names to their assigned phase.
type phaseRegistry struct {
	phases map[string]Phase
}

func newPhaseRegistry() *phaseRegistry {
	return &phaseRegistry{
		phases: make(map[string]Phase),
	}
}

// Assign associates a handler name with a shutdown phase.
func (r *phaseRegistry) Assign(handlerName string, phase Phase) {
	r.phases[handlerName] = phase
}

// Get returns the phase assigned to a handler, defaulting to PhaseApplication.
func (r *phaseRegistry) Get(handlerName string) Phase {
	if p, ok := r.phases[handlerName]; ok {
		return p
	}
	return PhaseApplication
}

// GroupByPhase partitions handler names into ordered phase buckets.
// Handlers within each bucket share the same phase and run concurrently;
// buckets themselves are executed sequentially in descending priority order.
func (r *phaseRegistry) GroupByPhase(names []string) [][]string {
	phaseBuckets := make(map[string][]string)
	phasePriority := make(map[string]int)

	for _, name := range names {
		p := r.Get(name)
		phaseBuckets[p.Name] = append(phaseBuckets[p.Name], name)
		phasePriority[p.Name] = p.Priority
	}

	// Collect unique phase names sorted by descending priority.
	type phaseMeta struct {
		name     string
		priority int
	}
	var metas []phaseMeta
	seen := make(map[string]bool)
	for pname, pri := range phasePriority {
		if !seen[pname] {
			metas = append(metas, phaseMeta{pname, pri})
			seen[pname] = true
		}
	}

	// Simple insertion sort — phase count is small.
	for i := 1; i < len(metas); i++ {
		for j := i; j > 0 && metas[j].priority > metas[j-1].priority; j-- {
			metas[j], metas[j-1] = metas[j-1], metas[j]
		}
	}

	result := make([][]string, 0, len(metas))
	for _, m := range metas {
		result = append(result, phaseBuckets[m.name])
	}
	return result
}

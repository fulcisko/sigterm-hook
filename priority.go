package sigterm

// Priority defines the shutdown priority level for a handler.
// Handlers with higher priority values are shut down first.
type Priority int

const (
	// PriorityLow is for non-critical background services.
	PriorityLow Priority = 0
	// PriorityNormal is the default priority for most handlers.
	PriorityNormal Priority = 50
	// PriorityHigh is for important services that should stop early.
	PriorityHigh Priority = 100
	// PriorityCritical is for handlers that must run first during shutdown.
	PriorityCritical Priority = 200
)

// priorityEntry associates a handler name with its priority.
type priorityEntry struct {
	name     string
	priority Priority
}

// sortByPriority performs a stable sort of handler names by descending priority.
// Handlers with equal priority retain their registration order.
func sortByPriority(names []string, priorities map[string]Priority) []string {
	if len(priorities) == 0 {
		return names
	}

	entries := make([]priorityEntry, len(names))
	for i, name := range names {
		p, ok := priorities[name]
		if !ok {
			p = PriorityNormal
		}
		entries[i] = priorityEntry{name: name, priority: p}
	}

	// Stable insertion sort — preserves registration order for equal priorities.
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].priority > entries[j-1].priority; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}

	sorted := make([]string, len(entries))
	for i, e := range entries {
		sorted[i] = e.name
	}
	return sorted
}

// String returns a human-readable label for the priority level.
// Named constants are returned as their label; other values are formatted
// as "Priority(<value>)".
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "Low"
	case PriorityNormal:
		return "Normal"
	case PriorityHigh:
		return "High"
	case PriorityCritical:
		return "Critical"
	default:
		return fmt.Sprintf("Priority(%d)", int(p))
	}
}

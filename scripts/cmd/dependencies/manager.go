package dependencies

type dependency interface {
	update() error
}

var (
	dependenciesList []dependency
)

// Manager is the type with knowledge about how to handle dependencies.
type Manager struct{}

// NewManager is the Manager constructor.
func NewManager() *Manager {
	return &Manager{}
}

// Run is the main entry point, it executes all the configured dependency
// updates.
func (*Manager) Run() error {
	for _, dep := range dependenciesList {
		if err := dep.update(); err != nil {
			return err
		}
	}
	return nil
}

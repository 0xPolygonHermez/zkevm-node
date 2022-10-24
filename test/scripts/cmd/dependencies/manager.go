package dependencies

type dependency interface {
	update() error
}

// Manager is the type with knowledge about how to handle dependencies.
type Manager struct {
	cfg *Config
}

// Config has the configurations options for all the updaters.
type Config struct {
	Images *ImagesConfig
	PB     *PBConfig
	TV     *TVConfig
}

// NewManager is the Manager constructor.
func NewManager(cfg *Config) *Manager {
	return &Manager{
		cfg: cfg,
	}
}

// Run is the main entry point, it executes all the configured dependency
// updates.
func (m *Manager) Run() error {
	iu := newImageUpdater(m.cfg.Images.Names, m.cfg.Images.TargetFilePath)
	pb := newPBUpdater(m.cfg.PB.SourceRepo, m.cfg.PB.TargetDirPath)
	tv := newTestVectorUpdater(m.cfg.TV.SourceRepo, m.cfg.TV.TargetDirPath)

	for _, dep := range []dependency{iu, pb, tv} {
		if err := dep.update(); err != nil {
			return err
		}
	}
	return nil
}

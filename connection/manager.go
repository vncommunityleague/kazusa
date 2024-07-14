package connection

type (
	managerDependencies interface {
		Repository
	}

	ManagementProvider interface {
		ConnectionManager() *Manager
	}

	Manager struct {
		d managerDependencies
	}
)

func NewManager(d managerDependencies) *Manager {
	return &Manager{d: d}
}

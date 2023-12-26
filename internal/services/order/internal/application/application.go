package application

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct{}

	Queries struct{}
)

func New() *Application {
	return &Application{
		Commands: Commands{},
	}
}

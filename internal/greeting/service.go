package greeting

type Service interface {
	GetGreeting() string
}

type serviceImpl struct{}

func New() Service {
	return &serviceImpl{}
}

func (svc *serviceImpl) GetGreeting() string {
	return "Greeting from agent!"
}

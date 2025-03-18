package quotes

type Service struct {
	deps Deps
}

type Deps interface {
	GetRandomQuote() string
}

func New(deps Deps) *Service {
	return &Service{
		deps: deps,
	}
}

func (s *Service) GetQuote() string {
	return s.deps.GetRandomQuote()

}

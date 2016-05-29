package soletta

type Soletta struct {
}

func NewSoletta() *Soletta {
	return &Soletta{}
}

func (s *Soletta) Start() bool {
	return true
}

func (s *Soletta) Stop() bool {
	return true
}

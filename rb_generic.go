package rb

type genericSrc struct {
}

func (s *genericSrc) Seed(seed int64) {
	panic("no seeding")
}

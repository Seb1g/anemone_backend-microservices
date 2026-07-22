// Package service contains business logic layers
package service

type service struct {
	CardRepo    CardRepo
	ColumnRepo  ColumnRepo
	BoardRepo   BoardRepo
	DefaultStep float64
}

func NewService(
	cardRepo CardRepo,
	columnRepo ColumnRepo,
	boardRepo BoardRepo,
	df float64,
) *service {
	return &service{
		CardRepo:    cardRepo,
		ColumnRepo:  columnRepo,
		BoardRepo:   boardRepo,
		DefaultStep: df,
	}
}

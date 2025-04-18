package usecase

type AccountRepository interface {
	Create()
	Update()
}

type AccountUsecase struct {
	repository AccountRepository
}

func (a *AccountUsecase) Create() error {
	a.repository.Create()
	return nil
}

func (a *AccountUsecase) Update() error {
	a.repository.Update()
	return nil
}

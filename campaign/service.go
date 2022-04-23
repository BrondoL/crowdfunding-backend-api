package campaign

type Service interface {
	GetCampaigns(id int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(id int) ([]Campaign, error) {
	var campaigns []Campaign
	var err error

	if id != 0 {
		campaigns, err = s.repository.FindByUserId(id)
	} else {
		campaigns, err = s.repository.FindAll()
	}

	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

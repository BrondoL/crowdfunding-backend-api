package campaign

import "errors"

type Service interface {
	GetCampaigns(id int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
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

func (s *service) GetCampaignById(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindById(input.ID)
	if err != nil {
		return campaign, err
	}
	if campaign.ID == 0 {
		return campaign, errors.New("user not found")
	}
	return campaign, nil
}

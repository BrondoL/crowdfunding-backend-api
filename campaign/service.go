package campaign

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(id int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
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

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.Perks = input.Perks
	campaign.UserID = input.User.ID

	timeNow := strconv.FormatInt(time.Now().Unix(), 10)
	slugName := fmt.Sprintf("%s %d %s", input.Name, input.User.ID, timeNow)
	campaign.Slug = slug.Make(slugName)

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

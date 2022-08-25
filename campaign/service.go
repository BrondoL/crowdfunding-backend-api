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
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	SaveImageCampaign(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
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

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindById(inputID.ID)
	if err != nil {
		return campaign, err
	}
	if campaign.ID == 0 {
		return campaign, errors.New("campaign not found")
	}

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("you are not allowed to update this campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.GoalAmount = inputData.GoalAmount
	campaign.Perks = inputData.Perks

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil
}

func (s *service) SaveImageCampaign(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	campaign, err := s.repository.FindById(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}
	if campaign.ID == 0 {
		return CampaignImage{}, errors.New("campaign not found")
	}
	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("you are not allowed to upload in this campaign")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		err := s.repository.MarkImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.FileName = fileLocation
	campaignImage.IsPrimary = isPrimary

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}
	return newCampaignImage, nil
}

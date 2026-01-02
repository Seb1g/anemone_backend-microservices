package services

import (
	"errors"
	"net/mail"
	"strings"

	"anemone_backend-microservices/internal/mail/model"
	"anemone_backend-microservices/internal/mail/repository"
)

type MailService struct {
	repo   *repository.Repository
	domain string
}

func New(repo *repository.Repository, domain string) *MailService {
	return &MailService{repo: repo, domain: domain}
}

func (s *MailService) CreateAddress(userID int64, address string) (*model.TempAddress, error) {
	parsed, err := mail.ParseAddress(address)
	if err != nil || parsed.Address != address {
		return nil, errors.New("invalid email format")
	}

	if !strings.HasSuffix(address, "@"+s.domain) {
		return nil, errors.New("invalid domain")
	}

	return s.repo.CreateAddress(address, userID)
}

func (s *MailService) GetInbox(addressID int) ([]model.Email, error) {
	return s.repo.GetEmailsForAddress(addressID)
}

func (s *MailService) ListAddresses(userID int64) ([]model.TempAddress, error) {
	return s.repo.GetAddressesForUser(userID)
}

func (s *MailService) DeleteAddress(addressID int, userID int64) error {
	ok, err := s.repo.CheckAddressOwner(addressID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("forbidden")
	}
	return s.repo.DeleteAddress(addressID, userID)
}

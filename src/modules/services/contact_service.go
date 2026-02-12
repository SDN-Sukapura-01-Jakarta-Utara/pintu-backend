package services

import (
	"encoding/json"
	"errors"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

type ContactService interface {
	Create(req *dtos.ContactCreateRequest, userID uint) (*dtos.ContactResponse, error)
	GetByID(id uint) (*dtos.ContactResponse, error)
	GetAll() ([]*dtos.ContactResponse, error)
	Update(id uint, req *dtos.ContactUpdateRequest, userID uint) (*dtos.ContactResponse, error)
	Delete(id uint) error
}

type ContactServiceImpl struct {
	repository repositories.ContactRepository
}

// NewContactService creates a new Contact service
func NewContactService(repository repositories.ContactRepository) ContactService {
	return &ContactServiceImpl{
		repository: repository,
	}
}

// Create creates a new Contact
func (s *ContactServiceImpl) Create(req *dtos.ContactCreateRequest, userID uint) (*dtos.ContactResponse, error) {
	// Convert jam_buka to JSON
	jamBukaJSON, _ := json.Marshal(req.JamBuka)

	// Create contact record
	data := &models.Contact{
		Alamat:    req.Alamat,
		Telepon:   req.Telepon,
		Email:     req.Email,
		JamBuka:   jamBukaJSON,
		Gmaps:     req.Gmaps,
		Website:   req.Website,
		Youtube:   req.Youtube,
		Instagram: req.Instagram,
		Tiktok:    req.Tiktok,
		Facebook:  req.Facebook,
		Twitter:   req.Twitter,
		CreatedByID: &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves Contact by ID
func (s *ContactServiceImpl) GetByID(id uint) (*dtos.ContactResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Contacts
func (s *ContactServiceImpl) GetAll() ([]*dtos.ContactResponse, error) {
	data, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]*dtos.ContactResponse, len(data))
	for i, item := range data {
		responses[i] = s.mapToResponse(&item)
	}

	return responses, nil
}

// Update updates Contact
func (s *ContactServiceImpl) Update(id uint, req *dtos.ContactUpdateRequest, userID uint) (*dtos.ContactResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("contact not found")
	}

	// Update basic fields if provided
	if req.Alamat != "" {
		existing.Alamat = req.Alamat
	}
	if req.Telepon != "" {
		existing.Telepon = req.Telepon
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if len(req.JamBuka) > 0 {
		jamBukaJSON, _ := json.Marshal(req.JamBuka)
		existing.JamBuka = jamBukaJSON
	}
	if req.Gmaps != "" {
		existing.Gmaps = req.Gmaps
	}
	if req.Website != "" {
		existing.Website = req.Website
	}
	if req.Youtube != "" {
		existing.Youtube = req.Youtube
	}
	if req.Instagram != "" {
		existing.Instagram = req.Instagram
	}
	if req.Tiktok != "" {
		existing.Tiktok = req.Tiktok
	}
	if req.Facebook != "" {
		existing.Facebook = req.Facebook
	}
	if req.Twitter != "" {
		existing.Twitter = req.Twitter
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes Contact by ID
func (s *ContactServiceImpl) Delete(id uint) error {
	// Get existing data
	_, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("contact not found")
	}

	// Delete from database
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *ContactServiceImpl) mapToResponse(data *models.Contact) *dtos.ContactResponse {
	// Map jam_buka from JSON
	var jamBukaItems []dtos.JamBukaItem
	if err := json.Unmarshal(data.JamBuka, &jamBukaItems); err != nil {
		jamBukaItems = []dtos.JamBukaItem{}
	}

	return &dtos.ContactResponse{
		ID:        data.ID,
		Alamat:    data.Alamat,
		Telepon:   data.Telepon,
		Email:     data.Email,
		JamBuka:   jamBukaItems,
		Gmaps:     data.Gmaps,
		Website:   data.Website,
		Youtube:   data.Youtube,
		Instagram: data.Instagram,
		Tiktok:    data.Tiktok,
		Facebook:  data.Facebook,
		Twitter:   data.Twitter,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		CreatedByID: data.CreatedByID,
		UpdatedByID: data.UpdatedByID,
	}
}

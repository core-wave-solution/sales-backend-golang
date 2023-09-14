package clientusecases

import (
	"context"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

type Service struct {
	rclient  cliententity.Repository
	rcontact personentity.ContactRepository
}

func NewService(rcliente cliententity.Repository, rcontact personentity.ContactRepository) *Service {
	return &Service{rclient: rcliente, rcontact: rcontact}
}

func (s *Service) RegisterClient(ctx context.Context, dto *clientdto.RegisterClientInput) (uuid.UUID, error) {
	client, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err := s.rclient.RegisterClient(ctx, client); err != nil {
		return uuid.Nil, err
	}

	for _, contact := range client.Contacts {
		if err := s.rcontact.RegisterContact(ctx, &contact); err != nil {
			return uuid.Nil, err
		}
	}

	return client.ID, nil
}

func (s *Service) UpdateClient(ctx context.Context, dtoId *entitydto.IdRequest, dto *clientdto.UpdateClientInput) error {
	client, err := s.rclient.GetClientById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(client); err != nil {
		return err
	}

	if err := s.rclient.UpdateClient(ctx, client); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteClient(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.rclient.GetClientById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rclient.DeleteClient(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetClientById(ctx context.Context, dto *entitydto.IdRequest) (*cliententity.Client, error) {
	if client, err := s.rclient.GetClientById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (s *Service) GetClientsBy(ctx context.Context, dto *clientdto.FilterClientInput) ([]cliententity.Client, error) {
	if clients, err := s.rclient.GetAllClients(ctx); err != nil {
		return nil, err
	} else {
		return clients, nil
	}
}

func (s *Service) GetAllClients(ctx context.Context) ([]cliententity.Client, error) {
	if clients, err := s.rclient.GetAllClients(ctx); err != nil {
		return nil, err
	} else {
		return clients, nil
	}
}

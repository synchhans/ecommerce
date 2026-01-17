package inventory

import "context"

type Repository interface {
	GetAvailability(ctx context.Context, variantID string) (*Availability, error)
}
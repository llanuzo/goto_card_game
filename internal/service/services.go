package service

import (
	"context"
)

type Services struct {
}

func NewServices() (s Services) {

	return s
}

func (s Services) Init(ctx context.Context) error {

	return nil
}

package cache

import (
	"context"
	"strconv"
)

type CreateSessionParams struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (cacher *RedisCacher) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	data := make(map[string]interface{})
	data["id"] = arg.ID
	data["email"] = arg.Email
	data["is_blocked"] = false

	if err := cacher.HMSet(ctx, arg.ID, data).Err(); err != nil {
		return Session{}, err
	}

	return Session{
		ID:        arg.ID,
		Email:     arg.Email,
		IsBlocked: false,
	}, nil
}

func (cacher *RedisCacher) GetSessionByID(ctx context.Context, id string) (Session, error) {
	var s Session

	data, err := cacher.HGetAll(ctx, id).Result()
	if err != nil {
		return s, err
	}

	s.ID = data["id"]
	s.Email = data["email"]
	s.IsBlocked, err = strconv.ParseBool(data["is_blocked"])
	if err != nil {
		return s, err
	}

	return s, nil
}

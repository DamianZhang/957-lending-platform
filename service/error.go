package service

import "errors"

// Service errors
var (
	ErrBadRequest      = errors.New("bad request")
	ErrInternalFailure = errors.New("internal failure")
)

type Error struct {
	svcErr error // service error
	appErr error // the reason of service error
}

func NewError(svcErr, appErr error) error {
	return Error{
		svcErr: svcErr,
		appErr: appErr,
	}
}

func (e Error) SvcErr() error {
	return e.svcErr
}

func (e Error) AppErr() error {
	return e.appErr
}

func (e Error) Error() string {
	return errors.Join(e.svcErr, e.appErr).Error()
}

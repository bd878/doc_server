package model

import "errors"

var (
	ErrWrongToken     = errors.New("wrong token")
	ErrNoUser         = errors.New("no user")
	ErrWrongPassword  = errors.New("wrong password")
)
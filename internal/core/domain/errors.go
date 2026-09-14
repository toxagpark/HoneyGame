package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserHoneyNotFound = errors.New("user honey not found")
	ErrNotEnoughHoney    = errors.New("not enough honey")

	ErrChallengeNotFound      = errors.New("challenge not found")
	ErrChallengeAlreadyExists = errors.New("challenge already exists")
	ErrChallengeNotOwner      = errors.New("not challenge creator")
	ErrSelfChallenge          = errors.New("cannot accept own challenge")
	ErrWrongAmount            = errors.New("wrong amount")
)

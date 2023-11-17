package authgrp

import (
	"fmt"
	"net/mail"

	"github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/web/v1/mid"
)

func toCoreNewUser(tokenInfo *auth.IDTokenInfo) (user.NewUser, error) {
	// upsert user and get user id here
	addr, err := mail.ParseAddress(tokenInfo.Email)
	if err != nil {
		return user.NewUser{}, mid.WrapError(fmt.Errorf("parsing email: %w", err))
	}
	usr := user.NewUser{
		Name:     tokenInfo.Name,
		Email:    *addr,
		ImageURL: tokenInfo.Picture,
		Sub:      tokenInfo.Sub,
	}

	return usr, nil
}

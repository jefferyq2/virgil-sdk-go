/*
 * Copyright (C) 2015-2018 Virgil Security Inc.
 *
 * Lead Maintainer: Virgil Security Inc. <support@virgilsecurity.com>
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 *   (1) Redistributions of source code must retain the above copyright
 *   notice, this list of conditions and the following disclaimer.
 *
 *   (2) Redistributions in binary form must reproduce the above copyright
 *   notice, this list of conditions and the following disclaimer in
 *   the documentation and/or other materials provided with the
 *   distribution.
 *
 *   (3) Neither the name of the copyright holder nor the names of its
 *   contributors may be used to endorse or promote products derived
 *   from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ''AS IS'' AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 */

package sdk

import (
	"gopkg.in/virgil.v5/errors"
	"time"
	"sync"
)

type CachingJwtProvider struct {
	RenewTokenCallback func(context *TokenContext) (string, error)
	Jwt *Jwt
	lock sync.Mutex
}

func NewCachingJwtProvider(renewTokenCallback func(context *TokenContext) (string, error)) *CachingJwtProvider {
	return &CachingJwtProvider{
		RenewTokenCallback:    renewTokenCallback,
	}
}

func (g *CachingJwtProvider) GetToken(context *TokenContext) (AccessToken, error) {

	if context == nil {
		return nil, errors.New("context is mandatory")
	}

	if g.RenewTokenCallback == nil {
		return nil, errors.New("callback is mandatory")
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.Jwt == nil || g.Jwt.IsExpiredDelta(5*time.Second) != nil {

		token, err := g.RenewTokenCallback(context)
		if err != nil{
			return nil, err

		}
		jwt, err := JwtFromString(token)
		if err != nil{
			return nil, err

		}
		g.Jwt = jwt
	}

	return g.Jwt, nil
}

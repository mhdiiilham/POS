package service

import (
	"github.com/mhdiiilham/POS/entity/user"
)

type apiService struct {
	userRepository user.Repository
	hasher         Hasher
	tokenSigner    TokenSigner
}

func NewAPIService(userRepository user.Repository, pwdHasher Hasher, tokenSigner TokenSigner) *apiService {
	return &apiService{
		userRepository: userRepository,
		hasher:         pwdHasher,
		tokenSigner:    tokenSigner,
	}
}

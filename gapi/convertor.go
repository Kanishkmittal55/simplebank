package gapi

import (
	db "github.com/kanishkmittal55/simplebank/db/sqlc"
	"github.com/kanishkmittal55/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangeAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

package gapi

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/kanishkmittal55/simplebank/pb"
	"github.com/kanishkmittal55/simplebank/val"
	"github.com/kanishkmittal55/simplebank/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) SendOtpMail(ctx context.Context, req *pb.SendOtpMailRequest) (*pb.SendOtpMailResponse, error) {
	violations := validateSendOtpMailParams(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	result, err := server.store.GetUserByEmail(ctx, req.GetEmailId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "This email is not registered with us: %s", req.GetEmailId())
	}

	taskPayload := &worker.PayloadSendOtpEmail{
		Username: result.Username,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err = server.taskDistributor.DistributeTaskSendOtpEmail(ctx, taskPayload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "User Found But Internal Server Error could not send otp email")
	}

	resp := pb.SendOtpMailResponse{
		MessageBody: fmt.Sprintf("Otp Email sent successfully"),
	}

	return &resp, nil
}

func validateSendOtpMailParams(req *pb.SendOtpMailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.EmailId); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}

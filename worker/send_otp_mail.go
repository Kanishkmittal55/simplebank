package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	db "github.com/kanishkmittal55/simplebank/db/sqlc"
	"github.com/kanishkmittal55/simplebank/db/util"
	"github.com/rs/zerolog/log"
)

const TaskSendOtpEmail = "task:send_otp_mail"

type PayloadSendOtpEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendOtpEmail(
	ctx context.Context,
	payload *PayloadSendOtpEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshalize payload : %w", err)
	}
	task := asynq.NewTask(TaskSendOtpEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task : %w", err)
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendOtpEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist: %w ", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	OtpEmail, err := processor.store.CreateOtpEmail(ctx, db.CreateOtpEmailParams{
		Username: user.Username,
		Email:    user.Email,
		OtpSent:  util.RandomOtp(3, 3),
	})
	if err != nil {
		return fmt.Errorf("failed to send Otp email: %w", err)
	}

	subject := "Welcome to HassleSkip"
	//verifyUrl := fmt.Sprintf("http://localhost:80/v1/otp_email?otp_id=%d&count=%s", OtpEmail.ID, OtpEmail.TotalOtpGenerations)
	content := fmt.Sprintf(`Hello %s, <br/>
	Thank You for your patience !<br/>
	Please enter this OTP to reset your password - %s <br/>`, user.FullName, OtpEmail.OtpSent)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email : %w", err)
	}
	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")

	return nil
}

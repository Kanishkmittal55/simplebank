-- name: CreateOtpEmail :one
INSERT INTO otp_emails (
    username,
    email,
    otp_sent
) VALUES (
     $1, $2, $3
 ) RETURNING *;
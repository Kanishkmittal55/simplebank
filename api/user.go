package api

import (
	"database/sql"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/kanishkmittal55/simplebank/worker"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/kanishkmittal55/simplebank/db/sqlc"
	"github.com/kanishkmittal55/simplebank/db/util"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

//func newUserResponse(user db.User, accessToken string) UserResponse {
//	return UserResponse{
//		AccessToken:      accessToken,
//		Username:         user.Username,
//		FullName:         user.FullName,
//		Email:            user.Email,
//		CreatedAt:        user.CreatedAt,
//		PasswordChangeAt: user.PasswordChangeAt,
//	}
//}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// First we calculate the hashed password using the util function
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := UserResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangedAt,
		CreatedAt:        user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)

}

type getUserRequest struct {
	Username string `uri:"usr" binding:"required,alphanum"`
}

func (server *Server) getUserByUsername(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// TODO : Implement the compelete Users API...!!! Modify , Delete, List All Users
type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken           string    `json:"access_token"`
	SessionID             uuid.UUID `json:"session_id"`
	User                  db.User   `json:"user"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// If everything is correct we will have to check the password provided by the client
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := LoginUserResponse{
		SessionID:             session.ID,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		AccessToken:           accessToken,
		User:                  user,
	}
	ctx.JSON(http.StatusOK, rsp)

}

type LogoutUserRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

type LogoutUserResponse struct {
	Message string `json:"message"`
	User    string `json:"user"`
}

func (server *Server) logoutUser(ctx *gin.Context) {
	var req LogoutUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Retrieve session details by session ID
	uid := uuid.MustParse(req.SessionID)
	session, err := server.store.GetSession(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Invalidate the session by setting is_blocked to true
	err = server.store.InvalidateSession(ctx, session.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := LogoutUserResponse{
		Message: "Successfully logged out from your session",
		User:    session.Username,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type ForgotPasswordRequest struct {
	EmailId string `json:"email_id" binding:"required"`
}

type ForgotPasswordResponse struct {
	MessageBody string `json:"message"`
}

func (server *Server) forgotPassword(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.EmailId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	taskPayload := &worker.PayloadSendOtpEmail{
		Username: user.Username,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err = server.taskDistributor.DistributeTaskSendOtpEmail(ctx, taskPayload, opts...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := ForgotPasswordResponse{
		MessageBody: fmt.Sprintf("Otp Email sent successfully"),
	}
	ctx.JSON(http.StatusOK, resp)
}

package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type getlandingPageRequest struct {
	key int64 `uri:"key" binding:"required,min=1"`
}

func (server *Server) getLandingPage(ctx *gin.Context) {
	var req getlandingPageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//arg := getlandingPageRequest{
	//	key: 1234567,
	//}

	message := "Hi User Welcome to API Bank"
	// account, err := server.store.CreateAccount(ctx, arg)

	ctx.JSON(http.StatusOK, message)

}

So below is the first implementation of this REST Api Backend Built in go - 

func (server *Server) createUser(ctx *gin.Context) {
var req createUserRequest
if err := ctx.ShouldBindJSON(&req); err != nil {
ctx.JSON(http.StatusBadRequest, errorResponse(err))
return
}

### Create User function without the use of db transaction to enqueue the job
	// First we calculate the hashed password using the util function
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Creating a User without using any transaction
	//arg := db.CreateUserParams{
	//	Username:       req.Username,
	//	HashedPassword: hashedPassword,
	//	FullName:       req.FullName,
	//	Email:          req.Email,
	//}
	//
	//user, err := server.store.CreateUser(ctx, arg)
	//if err != nil {
	//	if pqErr, ok := err.(*pq.Error); ok {
	//		switch pqErr.Code.Name() {
	//		case "unique_violation":
	//			ctx.JSON(http.StatusForbidden, errorResponse(err))
	//			return
	//		}
	//	}
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.Username,
			HashedPassword: hashedPassword,
			FullName:       req.FullName,
			Email:          req.Email,
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		status.Errorf(codes.Internal, "Failed to create User: %v", err)
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := UserResponse{
		Username:         txResult.User.Username,
		FullName:         txResult.User.FullName,
		Email:            txResult.User.Email,
		PasswordChangeAt: txResult.User.PasswordChangeAt,
		CreatedAt:        txResult.User.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)

}

## This was first changed to use the create_user_tx that was created to create a user and enqueue a send_verify_email_task together , so that one doesnot happen without the other. Now we had to make some integrations and changes in the main.go file also 

func runGinServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
    server, err := api.NewServer(config, store, taskDistributor)
    if err != nil {
        log.Fatal().Msg(("Cannot create server:"))
    }

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}

As you can see the worker.TaskDistributor module has been added to the gin server simply and then it has been used in the main() function with the goTaskProcessor instead of the using the GRPC server we took this approach -

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	// go runTaskProcessor(config, redisOpt, store)
	// go runGatewayServer(config, store, taskDistributor)
	// runGrpcServer(config, store, taskDistributor)

	// You can use runGinServer(config, store) instead of runGrpcServer(config, store) to serve http routed requests
	go runTaskProcessor(config, redisOpt, store)
	runGinServer(config, store, taskDistributor)
}

Now the first problem that we observed was different log formats , also the task is enqueued and processed and also the user is added successfully but over her we need to make some more edits...
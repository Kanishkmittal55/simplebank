package db

//
//import (
//	"context"
//	"database/sql"
//)
//
//// CreateUserTxParams contains the input parameters of the transfer transaction
//type SendOtpMailTxParams struct {
//	EmailId string
//}
//
//// CreateUserTxResult is the result of the transfer transaction
//type SendOtpMailTxResult struct {
//	user        User
//	MessageBody string
//}
//
//// CreateUserTx perform a money transfer from one account to the other.
//// It creates a transfer record, add account entries ( +, -) and update accounts' balance within a single database transaction
//func (store *SQLStore) SendOtpMailTx(ctx context.Context, arg SendOtpMailTxParams) (SendOtpMailTxResult, error) {
//	var result SendOtpMailTxResult
//
//	err := store.execTx(ctx, func(q *Queries) error {
//		var err error
//
//		result.user, err = q.GetUserByEmail(ctx, arg.EmailId)
//		if err != nil {
//			return err
//		}
//
//		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
//			Username: result.VerifyEmail.Username,
//			IsEmailVerified: sql.NullBool{
//				Bool:  true,
//				Valid: true,
//			},
//		})
//		return err
//	})
//	return result, err
//}

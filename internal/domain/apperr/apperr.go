// Package apperr contains errors in all app
package apperr

import "errors"

var (
	// Card errors
	ErrCardCreateFailed      = errors.New("card creation failed")
	ErrCardDeleteFailed      = errors.New("card delete failed")
	ErrCardNotFound          = errors.New("card not found")
	ErrCardRenameFailed      = errors.New("card rename failed")
	ErrCardMoveFailed        = errors.New("card move failed")
	ErrCardGetFailed         = errors.New("get card failed")
	ErrColumnNotFoundForCard = errors.New("column not found for card operation")

	// Columns errors
	ErrColumnCreateFailed = errors.New("column creation failed")
	ErrColumnDeleteFailed = errors.New("column delete failed")
	ErrColumnRenameFailed = errors.New("column rename failed")
	ErrColumnNotFound     = errors.New("column not found")
	ErrColumnMoveFailed   = errors.New("column move failed")
	ErrColumnGetFailed    = errors.New("get column failed")

	// Board errors
	ErrBoardNotFound     = errors.New("board not found")
	ErrBoardCreateFailed = errors.New("board creation failed")
	ErrBoardUpdateFailed = errors.New("board update failed")
	ErrBoardDeleteFailed = errors.New("board delete failed")
	ErrBoardRenameFailed = errors.New("board rename failed")
	ErrBoardGetFailed    = errors.New("get board failed")

	ErrBoardAccessDenied = errors.New("board not found or access denied")
	ErrInternal          = errors.New("something went wrong internally")
)

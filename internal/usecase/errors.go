package usecase

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidUserID      = errors.New("invalid user id")

	ErrFieldIsEmpty   = errors.New("name and address cannot be empty")
	ErrFieldIsTooLong = errors.New("name or address is too long")

	ErrInvalidTemplateName   = errors.New("invalid template name")
	ErrInvalidTemplatePrice  = errors.New("invalid template price")
	ErrInvalidTemplateID     = errors.New("invalid template id")
	ErrTemplateAlreadyExists = errors.New("template already exists")
	ErrTemplateNotFound      = errors.New("template not found")
	ErrFetchTemplates        = errors.New("failed to fetch templates")
	ErrCreateTemplate        = errors.New("failed to create template")
	ErrGetAllTemplate        = errors.New("failed to get all templates")

	ErrInvalidProductID       = errors.New("invalid product id")
	ErrInvalidProductQty      = errors.New("invalid product quantity")
	ErrCreateProduct          = errors.New("failed to create product")
	ErrFetchProducts          = errors.New("failed to fetch products")
	ErrProductNotFound        = errors.New("product not found")
	ErrProductAlreadyAssigned = errors.New("product already assigned")
	ErrDetachProducts         = errors.New("failed to detach products from rent point")

	ErrRentPointAlreadyExists = errors.New("rent point already exists")
	ErrRentPointNotFound      = errors.New("rent point not found")
	ErrCreateRentpoint        = errors.New("failed to create rent point")
	ErrFetchRentPoints        = errors.New("failed to fetch rent points")
	ErrInvalidRentpointID     = errors.New("invalid rentpoint id")
	ErrDeletionRentPoint      = errors.New("error deleting the rental point")

	ErrInvalidOrderID    = errors.New("invalid order id")
	ErrCreateOrder       = errors.New("failed to create order")
	ErrFetchOrders       = errors.New("failed to fetch orders")
	ErrOrderNotFound     = errors.New("order not found")
	ErrOrderNotActive    = errors.New("order not active")
	ErrInvalidOrderState = errors.New("invalid order state")
	ErrCompleteOrder     = errors.New("failed to complete order")
	ErrForbidden         = errors.New("access denied: you do not own this order")
)

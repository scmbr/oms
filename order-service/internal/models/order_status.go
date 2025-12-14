package models

type OrderStatus string

const (
	StatusCreated   OrderStatus = "CREATED"
	StatusReserved  OrderStatus = "RESERVED"
	StatusPaid      OrderStatus = "PAID"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusFailed    OrderStatus = "FAILED"
)

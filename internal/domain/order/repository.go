package orderentity

import "context"

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetOrderBy(ctx context.Context, o *Order) ([]Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
}

type DeliveryRepository interface {
	CreateDeliveryOrder(ctx context.Context, delivery *DeliveryOrder) error
	UpdateDeliveryOrder(ctx context.Context, delivery *DeliveryOrder) error
	DeleteDeliveryOrder(ctx context.Context, id string) error
	GetDeliveryById(ctx context.Context, id string) (*DeliveryOrder, error)
	GetAllDeliveries(ctx context.Context) ([]DeliveryOrder, error)
}

type TableRepository interface {
	CreateTableOrder(ctx context.Context, table *TableOrder) error
	UpdateTableOrder(ctx context.Context, table *TableOrder) error
	DeleteTableOrder(ctx context.Context, id string) error
	GetTableById(ctx context.Context, id string) (*TableOrder, error)
	GetAllTables(ctx context.Context) ([]TableOrder, error)
}

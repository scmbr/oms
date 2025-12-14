package rabbit

const (
	orderQueue   = "order.service"
	sagaExchange = "saga.events"
)

var routingKeys = []string{
	"inventory.reserved",
	"inventory.failed",
	"payment.succeeded",
	"payment.failed",
}

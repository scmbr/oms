package worker

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/scmbr/oms/order-service/internal/service"
	"github.com/scmbr/oms/order-service/internal/transport/rabbit"
)

type OutboxWorker struct {
	outboxService service.Outbox
	publisher     *rabbit.Publisher
	interval      time.Duration
}

func NewOutboxWorker(outboxService service.Outbox, publisher *rabbit.Publisher, interval time.Duration) *OutboxWorker {
	return &OutboxWorker{
		outboxService: outboxService,
		publisher:     publisher,
		interval:      interval,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Outbox worker stopped")
			return
		case <-ticker.C:
			w.processPendingEvents(ctx)
		}
	}
}

func (w *OutboxWorker) processPendingEvents(ctx context.Context) {
	events, err := w.outboxService.GetPending(ctx)
	if err != nil {
		log.Println("Failed to fetch pending events:", err)
		return
	}

	for _, event := range events {
		if err := w.publisher.Publish(event.Payload, strconv.FormatUint(uint64(event.ID), 10)); err != nil {
			log.Println("Failed to publish event", event.ID, err)
			continue
		}

		if err := w.outboxService.MarkAsProcessed(ctx, nil, strconv.FormatUint(uint64(event.ID), 10)); err != nil {
			log.Println("Failed to mark event as processed", event.ID, err)
		}
	}
}

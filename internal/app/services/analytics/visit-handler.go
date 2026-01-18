package analytics

import (
	"context"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

type LinkVisitedProcessor struct {
	visitor LinkVisitor
}

func NewLinkVisitedProcessor(visitor LinkVisitor) *LinkVisitedProcessor {
	return &LinkVisitedProcessor{visitor: visitor}
}

func (processor LinkVisitedProcessor) Handler(event domain.LinkVisited) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	short_code, err := domain.NewShortCodeFromParam(event.ShortCode)

	if err != nil {
		return err
	}

	err = processor.visitor.Visit(ctx, short_code, event.VisitedAt)

	if err != nil {
		return err
	}

	return nil
}

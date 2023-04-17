package parcour

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ThinkChaos/parcour/jobgroup"
)

// Unbuffered can be used with `NewProducersWithBuffer` to obtain
// an unbuffered `Producers`.
const Unbuffered = 0

// Producer is a function that produces a series of `T`.
type Producer[T any] func(context.Context, chan<- T) error

// Consumer is a function that consumes a series of `T`.
type Consumer[T any] func(context.Context, <-chan T) error

// Producers implements a multiple producer, multiple consumer pattern.
//
// Producers create a series of items that the consumers... consume!
type Producers[T any] struct {
	producersGrp jobgroup.JobGroup
	consumersGrp jobgroup.JobGroup

	items chan T
	close sync.Once
}

// NewUnbufferedProducers returns a new `Producers`.
//
// The producers are unbuffered, meaning they will block after producing each item
// until a consumer receives it.
//
// This function is equivalent to `NewProducersWithBuffer(parentGroup, Unbuffered)`.
func NewUnbufferedProducers[T any](producersGrp, consumersGrp jobgroup.JobGroup) *Producers[T] {
	return NewProducersWithBuffer[T](producersGrp, consumersGrp, Unbuffered)
}

// NewProducersWithBuffer returns a new `Producers` with the given buffer capacity.
//
// The buffer capacity configures backpressure. It is the maximum number of items
// produced without a consumer receiving them.
// Once that limit is reached, producers will block on send, waiting for a consumer to receive.
//
// This function panics if `bufferCap` is negative.
func NewProducersWithBuffer[T any](producersGrp, consumersGrp jobgroup.JobGroup, bufferCap int) *Producers[T] {
	return &Producers[T]{
		producersGrp: jobgroup.WithParent(producersGrp),
		consumersGrp: jobgroup.WithParent(consumersGrp),

		items: make(chan T, bufferCap), // panics if bufferCap is negative
		close: sync.Once{},
	}
}

// Close waits for the producers, closes the items channel, and then waits for consumers.
func (p *Producers[T]) Close() {
	defer func() {
		// Broadcast end to consumers
		p.closeItems()

		p.consumersGrp.Close()
	}()

	p.producersGrp.Close()
}

func (p *Producers[T]) closeItems() {
	p.close.Do(func() {
		close(p.items)
	})
}

// BufferCap returns the receiver's buffer capacity.
//
// The result is undefined if the receiver's `Wait` or `Close` methods were called,
// though it is guaranteed the function will not panic in that case.
func (p *Producers[T]) BufferCap() int {
	return cap(p.items)
}

// GoProduce starts a new producer job.
func (p *Producers[T]) GoProduce(producer Producer[T]) {
	p.producersGrp.Go(func(ctx context.Context) error {
		return producer(ctx, p.items)
	})
}

// GoConsume starts a new consumer job.
func (p *Producers[T]) GoConsume(consumer Consumer[T]) {
	p.consumersGrp.Go(func(ctx context.Context) error {
		return consumer(ctx, p.items)
	})
}

// Wait blocks the current goroutine until all producer and consumer jobs to finish.
func (p *Producers[T]) Wait() error {
	err := wrapProducersError(p.producersGrp.Wait())

	// Broadcast end to consumers
	p.closeItems()

	return errors.Join(err, wrapConsumersError(p.consumersGrp.Wait()))
}

// ProducersError holds any errors returned by producer goroutines.
type ProducersError struct {
	inner error
}

func wrapProducersError(err error) error {
	if err == nil {
		return nil
	}

	return &ProducersError{inner: err}
}

func (e *ProducersError) Error() string {
	return fmt.Sprintf("producer error(s): %s", e.inner)
}

func (e *ProducersError) Unwrap() error {
	return e.inner
}

// ConsumersError holds any errors returned by consumer goroutines.
type ConsumersError struct {
	inner error
}

func wrapConsumersError(err error) error {
	if err == nil {
		return nil
	}

	return &ConsumersError{inner: err}
}

func (e *ConsumersError) Error() string {
	return fmt.Sprintf("consumer error(s): %s", e.inner)
}

func (e *ConsumersError) Unwrap() error {
	return e.inner
}

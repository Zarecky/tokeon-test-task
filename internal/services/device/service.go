package device

import (
	"context"
	"sync"
	"sync/atomic"
	"tokeon-test-task/internal/errors"

	"github.com/google/uuid"
)

type channel struct {
	message chan string
	stop    chan struct{}
}

type Service struct {
	devicesChannels map[uuid.UUID]channel
	mu              sync.RWMutex
}

func New() *Service {
	return &Service{
		devicesChannels: make(map[uuid.UUID]channel),
		mu:              sync.RWMutex{},
	}
}

func (s *Service) Register(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.devicesChannels[id]
	if ok {
		return errors.ErrDeviceAlreadyRegistered
	}

	s.devicesChannels[id] = channel{
		make(chan string),
		make(chan struct{}, 1),
	}

	return nil
}

func (s *Service) Get(id uuid.UUID) (<-chan string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ch, ok := s.devicesChannels[id]
	if !ok {
		return nil, errors.ErrDeviceNotFound
	}

	return ch.message, nil
}

func (s *Service) Close(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch, ok := s.devicesChannels[id]
	if !ok {
		return errors.ErrDeviceNotFound
	}

	ch.stop <- struct{}{}
	close(ch.stop)
	delete(s.devicesChannels, id)

	return nil
}

func (s *Service) SendMessage(ctx context.Context, deviceID *uuid.UUID, text string) error {
	var channels []channel

	if deviceID != nil {
		s.mu.RLock()

		ch, ok := s.devicesChannels[*deviceID]
		if !ok {
			return errors.ErrDeviceNotFound
		}

		s.mu.RUnlock()

		channels = []channel{ch}
	} else {
		s.mu.RLock()

		channels = make([]channel, 0, len(s.devicesChannels))
		for _, ch := range s.devicesChannels {
			channels = append(channels, ch)
		}

		s.mu.RUnlock()
	}

	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	var lastErr error
	var errCount int32

	for _, ch := range channels {

		go func(ctx context.Context, channel channel) {
			defer wg.Done()

			if err := s.send(ctx, channel, text); err != nil {
				lastErr = err
				atomic.AddInt32(&errCount, 1)
			}
		}(ctx, ch)
	}

	wg.Wait()

	if lastErr != nil && len(channels) == int(errCount) {
		return lastErr
	}

	return nil
}

func (s *Service) send(ctx context.Context, channel channel, text string) error {
	select {
	case channel.message <- text:
		return nil
	case <-channel.stop:
		return errors.ErrDeviceNotFound
	case <-ctx.Done():
		return nil
	}
}

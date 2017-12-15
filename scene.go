package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	paused bool
	life   *lifePainter
}

func newScene(w, h int) *scene {
	l := NewLife(w, h)
	return &scene{
		paused: true,
		life:   &lifePainter{l},
	}
}

func (s *scene) handleEvents() error {
	for {
		event := sdl.PollEvent()
		if event == nil {
			break
		}

		switch event.(type) {
		case *sdl.QuitEvent:
			return fmt.Errorf("quit")
		case *sdl.KeyboardEvent:
			k := event.(*sdl.KeyboardEvent)
			if k.Keysym.Sym == sdl.K_SPACE && k.Type == sdl.KEYDOWN && k.Repeat == 0 {
				s.paused = !s.paused
			}
		}
	}

	return nil
}

func (s *scene) run(r *sdl.Renderer) error {
	if err := s.paint(r); err != nil {
		return err
	}

	for {
		if err := s.handleEvents(); err != nil {
			return err
		}

		sdl.Delay(1000 / 60)

		if s.paused {
			continue
		}

		if err := s.paint(r); err != nil {
			return err
		}

		s.life.Step()
	}
}

func (s *scene) paint(r *sdl.Renderer) error {
	if err := s.life.paint(r); err != nil {
		return fmt.Errorf("could not paint lift: %s", err)
	}

	r.Present()
	return nil
}

type lifePainter struct {
	*Life
}

func (lr *lifePainter) paint(r *sdl.Renderer) error {
	walk := func(x, y int, alive bool) error {
		if !alive {
			return nil
		}

		rect := &sdl.Rect{
			int32(x * blockSize),
			int32(y * blockSize),
			blockSize,
			blockSize,
		}
		return r.FillRect(rect)
	}

	err := r.SetDrawColor(255, 255, 255, 255)
	if err != nil {
		return fmt.Errorf("could not set draw color: %v", err)
	}
	err = r.Clear()
	if err != nil {
		return fmt.Errorf("could not clear: %v", err)
	}

	err = r.SetDrawColor(0, 0, 0, 255)
	if err != nil {
		return fmt.Errorf("could not set draw color: %v", err)
	}
	err = lr.Life.Walk(walk)
	if err != nil {
		return fmt.Errorf("could not walk life: %v", err)
	}
	return nil
}

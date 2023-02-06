package main

import (
	"sync"
	"testing"
)

var (
	mu     sync.Mutex
	rwMu   sync.RWMutex
	number int
)

func write() {
  mu.Lock()
  defer mu.Unlock()

  number++
}

func read(){
  mu.Lock()
  defer mu.Unlock()

  _ = number
}

func rwRead()  {
  rwMu.RLock()
  defer rwMu.RUnlock()

  _ = number
}

func BenchmarkMutex10Write90Read(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%10 == 0 {
      write()
		} else {
			read()
		}
	}
}

func BenchmarkRWMutex10Write90Read(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%10 == 0 {
      write()
		} else {
			rwRead()
		}
	}
}

func BenchmarkMutex50Write50Read(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
      write()
		} else {
			read()
		}
	}
}

func BenchmarkRWMutex50Write50Read(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
      write()
		} else {
			rwRead()
		}
	}
}

func BenchmarkMutex90Write10Read(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%10 != 0 {
      write()
		} else {
			read()
		}
	}
}

func BenchmarkRWMutex90Write10Read(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%10 != 0 {
      write()
		} else {
			rwRead()
		}
	}
}

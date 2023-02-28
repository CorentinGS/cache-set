package cacheset

import (
	"reflect"
	"testing"
	"time"
)

func Test_newSet(t *testing.T) {
	type testCase[T comparable] struct {
		want set[T]
		name string
	}
	tests := []testCase[int64]{
		{
			name: "test",
			want: make(map[int64]int64),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSet[int64](); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Add(t *testing.T) {
	type args[T comparable] struct {
		elem     T
		duration time.Duration
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
	}
	s := newSet[int64]()

	tests := []testCase[int64]{
		{
			name: "test",
			args: args[int64]{
				elem:     1,
				duration: 0,
			},
		},
		{
			name: "test",
			args: args[int64]{
				elem:     1,
				duration: 1 * time.Minute,
			},
		},
		{
			name: "test",
			args: args[int64]{
				elem:     3,
				duration: 1 * time.Minute,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Add(tt.args.elem, tt.args.duration)

			if !s.Contains(tt.args.elem) {
				t.Errorf("Add() = %v, want %v", s, tt.args.elem)
			}
		})
	}
}

func Test_set_Clear(t *testing.T) {
	s := newSet[int64]()
	s.Add(1, 0)
	s.Add(2, 0)

	t.Run("Clear", func(t *testing.T) {
		s.Clear()
		if len(s) != 0 {
			t.Errorf("Clear() = %v, want %v", s, 0)
		}
	})
}

func Test_set_Contains(t *testing.T) {
	s := newSet[int64]()
	s.Add(1, 0)
	s.Add(2, 0)

	t.Run("Contains", func(t *testing.T) {
		if got := s.Contains(1); !got {
			t.Errorf("Contains() = %v, want %v", got, true)
		}
	})

	t.Run("Contains", func(t *testing.T) {
		if got := s.Contains(3); got {
			t.Errorf("Contains() = %v, want %v", got, false)
		}
	})
}

func Test_set_Copy(t *testing.T) {
	s := newSet[int64]()
	s.Add(1, 0)
	s.Add(2, 0)

	t.Run("Copy", func(t *testing.T) {
		if got := s.Copy(); !reflect.DeepEqual(got, s) {
			t.Errorf("Copy() = %v, want %v", got, s)
		}
	})
}

func Test_set_Delete(t *testing.T) {
	s := newSet[int64]()
	s.Add(1, 0)
	s.Add(2, 0)

	t.Run("Delete", func(t *testing.T) {
		s.Delete(1)
		if s.Contains(1) {
			t.Errorf("Delete() = %v, want %v", s, 0)
		}
	})
}

func Test_set_Expire(t *testing.T) {
	t.Parallel()
	s := newSet[int64]()
	s.Add(1, 1*time.Second)
	s.Add(2, 1*time.Minute)

	time.Sleep(2 * time.Second)

	s.Expire(1)
	s.Expire(2)

	t.Run("Expire", func(t *testing.T) {
		if s.Contains(1) {
			t.Errorf("Expire() = %v, want %v", s, 0)
		}
	})

	t.Run("Expire", func(t *testing.T) {
		if !s.Contains(2) {
			t.Errorf("Expire() = %v, want %v", s, 0)
		}
	})
}

func Test_set_ExpireAll(t *testing.T) {
	t.Parallel()
	s := newSet[int64]()
	s.Add(1, 1*time.Second)
	s.Add(2, 1*time.Minute)

	time.Sleep(2 * time.Second)

	s.ExpireAll()

	t.Run("ExpireAll", func(t *testing.T) {
		if s.Contains(1) {
			t.Errorf("ExpireAll() = %v, want %v", s, 0)
		}
	})

	t.Run("ExpireAll", func(t *testing.T) {
		if !s.Contains(2) {
			t.Errorf("ExpireAll() = %v, want %v", s, 0)
		}
	})
}

func Test_set_Expired(t *testing.T) {
	t.Parallel()
	s := newSet[int64]()
	s.Add(1, 1*time.Second)
	s.Add(2, 1*time.Minute)

	time.Sleep(2 * time.Second)

	t.Run("Expired", func(t *testing.T) {
		if !s.Expired(1) {
			t.Errorf("Expired() = %v, want %v", s, 0)
		}
	})

	t.Run("Expired", func(t *testing.T) {
		if s.Expired(2) {
			t.Errorf("Expired() = %v, want %v", s, 0)
		}
	})
}

func Test_set_Len(t *testing.T) {
	s := newSet[int64]()
	s.Add(1, 0)
	s.Add(2, 0)

	t.Run("Len", func(t *testing.T) {
		if got := s.Len(); got != 2 {
			t.Errorf("Len() = %v, want %v", got, 2)
		}
	})
}

func Test_set_ToSlice(t *testing.T) {
	s := newSet[int64]()
	s.Add(1, 0)
	s.Add(2, 0)

	t.Run("ToSlice", func(t *testing.T) {
		if got := s.ToSlice(); !reflect.DeepEqual(got, []int64{1, 2}) {
			t.Errorf("ToSlice() = %v, want %v", got, []int64{1, 2})
		}
	})
}

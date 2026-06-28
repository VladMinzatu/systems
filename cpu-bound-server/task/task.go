package task

import (
	"fmt"
	"math/rand/v2"
	"strconv"
)

type TaskProvider struct{}

func NewTaskProvider() *TaskProvider {
	return &TaskProvider{}
}

func (t *TaskProvider) GetTask(kind string, size int) (Task, error) {
	switch kind {
	case "cpu":
		return NewMatMulTask(size), nil
	case "sprintf":
		return NewSprintfTask(size), nil
	default:
		return nil, fmt.Errorf("Unrecognized task type requested")
	}
}

type Task interface {
	Run()
	Result() <-chan error
}

type MatMulTask struct {
	N int

	A      []float64
	B      []float64
	C      []float64
	result chan error
}

func NewMatMulTask(n int) *MatMulTask {
	task := MatMulTask{
		N:      n,
		A:      make([]float64, n*n),
		B:      make([]float64, n*n),
		C:      make([]float64, n*n),
		result: make(chan error, 1),
	}
	return &task
}

func (t *MatMulTask) Run() {
	n := t.N
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum +=
					t.A[i*n+k] *
						t.B[k*n+j]
			}
			t.C[i*n+j] = sum
		}
	}
	t.result <- nil
}

func (t *MatMulTask) Result() <-chan error {
	return t.result
}

type SprintfTask struct {
	Iterations int
	result     chan error
}

func NewSprintfTask(iterations int) *SprintfTask {
	return &SprintfTask{Iterations: iterations, result: make(chan error, 1)}
}

func (t *SprintfTask) Run() {
	for i := 0; i < t.Iterations; i++ {
		fmt.Sprintf(
			"%d-%d-%d-%f-%s",
			i,
			rand.Int(),
			rand.Int(),
			rand.Float64(),
			strconv.Itoa(i),
		)
	}
	t.result <- nil
}

func (t *SprintfTask) Result() <-chan error {
	return t.result
}

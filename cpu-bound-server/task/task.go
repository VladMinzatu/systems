package task

import "context"

type Task interface {
	Execute(ctx context.Context)
}

type MatMulTask struct {
	N int

	A []float64
	B []float64
	C []float64
}

func NewMatMulTask(n int) *MatMulTask {
	task := MatMulTask{
		N: n,
		A: make([]float64, n*n),
		B: make([]float64, n*n),
		C: make([]float64, n*n),
	}
	return &task
}

func (t *MatMulTask) Execute(ctx context.Context) error {
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

	return nil
}

package util

import "golang.org/x/sync/errgroup"

func Batch(inputs []string, f func(string) error) error {
	var g errgroup.Group
	for _, input := range inputs {
		g.Go(func() error {
			return f(input)
		})
	}
	return g.Wait()
}

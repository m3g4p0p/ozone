package cleanup

import (
	"errors"
	"io"
)

type Cleanup []io.Closer

func (c Cleanup) Close() error {
	var errs []error
	for _, c := range c {
		errs = append(errs, c.Close())
	}
	return errors.Join()
}

func (c *Cleanup) Add(closer io.Closer) {
	*c = append(*c, closer)
}

func (c *Cleanup) Take() io.Closer {
	closer := *c
	*c = nil
	return closer
}

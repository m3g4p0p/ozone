package resources

import "io"

func CloseAny(v any) error {
	if c, ok := v.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

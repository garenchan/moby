package ratelimit

import "io"

type reader struct {
	r       io.Reader
	limiter *Limiter
}

// Reader returns a reader that is rate limited.
// Each token in the limiter bucket represents one byte.
func Reader(r io.Reader, limiter *Limiter) io.ReadCloser {
	return &reader{r: r, limiter: limiter}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}

	r.limiter.Wait(int64(n))
	return n, err
}

func (r *reader) Close() error {
	rc, ok := r.r.(io.ReadCloser)
	if ok {
		return rc.Close()
	}

	return nil
}

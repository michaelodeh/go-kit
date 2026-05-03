package llmstream

import (
	"context"
	"io"
)

func (s *Client) StreamWithContext(
	ctx context.Context,
	in io.Reader,
	out io.Writer,
) error {
	buf := make([]byte, 32*1024)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := in.Read(buf)

		if n > 0 {
			if _, wErr := out.Write(buf[:n]); wErr != nil {
				return wErr
			}
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (s *Client) Stream(in io.Reader, out io.Writer) error {
	return s.StreamWithContext(context.Background(), in, out)
}

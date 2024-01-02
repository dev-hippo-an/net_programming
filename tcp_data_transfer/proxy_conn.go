package tcp_data_transfer

import (
	"io"
	"net"
)

func proxyConn(source, destination string) error {
	connSource, err := net.Dial("tcp", source)

	if err != nil {
		return err
	}

	defer connSource.Close()

	connDestination, err := net.Dial("tcp", destination)

	if err != nil {
		return err
	}

	defer connDestination.Close()

	// destination 에서 읽고 source 로 쓴다.
	go func() {
		_, _ = io.Copy(connSource, connDestination)
	}()

	// source 에서 읽고 destination 으로 쓴다.
	_, err = io.Copy(connDestination, connSource)

	return err
}

func proxy(from io.Reader, to io.Writer) error {
	fromWriter, fromIsWriter := from.(io.Writer)

	toReader, toIsReader := to.(io.Reader)

	if toIsReader && fromIsWriter {
		go func() {
			_, _ = io.Copy(fromWriter, toReader)
		}()
	}

	_, err := io.Copy(to, from)

	return err
}

/*

	Very simple program for automatically spinning up a h.264 stream
	when requested on port 10001. Kills the stream when the client
	disconnects, saving CPU time / power.

	Alter the values in the Command to suit. The location of the
	command is probably Arch specific, might need changing on
	another distro. To compile on Mac or Linux:

	GOOS=linux GOARM=6 GOARCH=arm go build <filename_of_this_file>

	Then copy the file to the Raspberry Pi. When you want to stream
	video run:

	nc <IP address of pi> 10001 | mplayer -fs -fps 200 -demuxer h264es -fs -

*/

package v_stream

import (
	"io"
	"log"
	"net"
	"os/exec"
	"strconv"
	"syscall"
)

var read_chunk int64 = 500

func ServeVideo(x_res, y_res, bitrate int) {
	log.Println("Starting video stream server")

	tcp_addr, err := net.ResolveTCPAddr("tcp4", ":10001")
	if err != nil {
		log.Fatal("Couldn't resolve tcp address, ", err)
	}

	ln, err := net.ListenTCP("tcp4", tcp_addr)

	// ln, err := net.Listen("tcp", ":10001")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Fatal("Error accepting connection, ", err)
		}

		conn.SetNoDelay(true)
		conn.SetWriteBuffer(10e7)

		cmd := exec.Command("/opt/vc/bin/raspivid", "-g", "10", "-n", "-w", strconv.Itoa(x_res), "-h", strconv.Itoa(y_res), "-b", strconv.Itoa(bitrate), "-fps", "30", "-t", "0", "-o", "-")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal("raspivid failed to start: ", err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		written, err := io.Copy(conn, stdout)
		if err != nil {
			log.Println("User disconnected, ", err)
		}

		log.Printf("Wrote %v bytes", written)

		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Fatal(err)
		}

		cmd.Process.Wait()
	}
}

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

package main

import (
	"io"
	"log"
	"net"
	"os/exec"
	"syscall"
)

func main() {
	ln, err := net.Listen("tcp", ":10001")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command("/opt/vc/bin/raspivid", "-n", "--hflip", "--vflip", "-w", "1280", "-h", "720", "-b", "1000000", "-t", "0", "-o", "-")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
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

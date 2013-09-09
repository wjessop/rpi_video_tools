# Raspberry Pi Camera module tools

A colection of programs for doing stuff with the Raspberry Pi camera module from Go.

Currently only one program, there may be more added at some point if I need to write them (or someone sends a pull request).

## Building the tools

	GOOS=linux GOARM=6 GOARCH=arm go build <filename_of_the_file>

then copy the binaries to the Pi.

## The tools

### v_stream

Very simple program for automatically spinning up a h.264 stream
when requested on port 10001. Kills the stream when the client
disconnects, saving CPU time / power.

Alter the values in the Command to suit. The location of the
command is probably Arch specific, might need changing on
another distro. To compile on Mac or Linux:

Run the program on the Pi, then to get a stream on your local machine run:

	nc <IP address of pi> 10001 | mplayer -fs -fps 200 -demuxer h264es -fs -

Only one stream supported at a time (deliberately). All other connections will wait.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## Author

* Will Jessop, @will_j, will@willj.net

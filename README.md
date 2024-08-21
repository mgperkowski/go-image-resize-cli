# go-image-resize-cli

A command line tool to resize images in bulk.

This tool does not modify the original images, it simply creates a new directory inside the directory it was called from called "Resized_Images" and places the resized images there.

## Installation

**Installation requires Go and Git to be installed on the host computer.

Clone the repository:

```bash
git clone https://github.com/mgperkowski/go-image-resize-cli.git
```
Navigate to the project directory:

```bash
cd go-image-resize-cli
```
Build the binary:

```bash
go build -o go-resize .
```

Move the binary to a directory in your PATH:

```bash
mv go-resize /usr/local/bin/
```

Run the CLI:

```bash
go-resize ./path/to/dir -w 500

#or

go-resize ./path/to/image.jpg -h 450
```

You can specify either the width or the height, the image aspect ratio will be maintained.  You can also give a path to a directory or a single image file.

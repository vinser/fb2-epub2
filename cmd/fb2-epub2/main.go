package main

import (
	"flag"
	"io"
	"os"

	"log"

	"github.com/vinser/fb2-epub2/pkg/fb2"
)

func main() {
	var (
		src string
		dst string
	)
	flag.StringVar(&src, "f", "", `fb2 file name`)
	flag.StringVar(&dst, "e", "", `epub file name`)
	flag.Parse()

	r, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	if err := ConvertFb2Epub(r, w); err != nil {
		log.Fatal(err)
	}
}

func ConvertFb2Epub(r io.ReadSeekCloser, w io.WriteCloser) error {
	fb, err := fb2.New(r)
	if err != nil {
		log.Fatal(err)
	}

	if err := fb.MakeEpub(w); err != nil {
		log.Fatal(err)
	}
	return nil
}

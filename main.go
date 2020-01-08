package main

import (
	"log"
	"os"

	"github.com/gproessl/torrent-client/torrent"
)

func printUsage() {
	log.Println("usage: torrent-client torrentfile outfile")
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}

	inPath := os.Args[1]
	outPath := os.Args[2]

	tf, err := torrent.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

	err = tf.DownloadToFile(outPath)
	if err != nil {
		log.Fatal(err)
	}
}

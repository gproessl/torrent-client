package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"golang.org/x/net/proxy"

	"github.com/jackpal/bencode-go"
)

const Port = 6887

// Torrent encodes the metadata from a .torrent file
type Torrent struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// DownloadToFile downloads a torrent and writes it to a file
func (t *Torrent) DownloadToFile(path string) error {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, nil)
	if err != nil {
		return err
	}

	buf, err := t.Download(dialer)
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// Open parses a torrent file
func Open(path string) (Torrent, error) {
	file, err := os.Open(path)
	if err != nil {
		return Torrent{}, err
	}
	defer file.Close()

	bto := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return Torrent{}, err
	}
	return bto.toTorrent()
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (bto *bencodeTorrent) toTorrent() (Torrent, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return Torrent{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return Torrent{}, err
	}
	t := Torrent{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}

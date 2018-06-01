package main

import (
	"bytes"
	"encoding/binary"
	"bufio"
	"os"
	"fmt"
	"strconv"
	"io/ioutil"
)

func readUint32(data []byte) (ret uint32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readUint16(data []byte) (ret uint16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

type header struct {
	sectionType	uint32
	chunkSize	uint32
	build		uint32
}

type txdFile struct {
	Header		header
	Info		txdInfo
	Textures	[]txdTexture
	Extra		txdExtraInfo
}

type txdInfo struct {
	Header		header
	Count		uint16
	Unknown		uint16
}

type txdTexture struct {
	Header		header
	Data		txdTextureData
	Extra		txdExtraInfo
}

type txdTextureData struct {
	Header			header
	Version			uint32
	FilterFlags		uint32
	TextureName		string
	AlphaName		string
	AlphaFlags		uint32
	TextureFormat	string
	Width			uint16
	Height			uint16
	Depth			uint8
	MipmapCount		uint8
	TexcodeType		uint8
	Flags			uint8
	Palette			uint8
	DataSize		uint32
	_DataStart		uint32
	_DataEnd		uint32
	/* Texture has to be stored here */
}

type txdExtraInfo struct {
	Header		header
}

func (h *header) read(p *bufio.Reader) bool {
	for {
		buf := make([]byte, 12)
		numBytesRead, err := p.Read(buf)
		if numBytesRead != 12 {
			return false
		}
		check(err)
		filePosition += uint32(numBytesRead)

		h.sectionType = readUint32(buf[:4])
		h.chunkSize = readUint32(buf[4:8])
		h.build = readUint32(buf[8:])

		//fmt.Printf("Header red! %+v\n", h)

		if h.sectionType != 3 {
			break
		}
	}
	return true
}

func (h *txdFile) read(f *os.File) bool {
	filePosition = 0
	p := bufio.NewReader(f)
	h.Header.read(p)
	if debug {
		fmt.Println("------------------------------------------")
	}
	h.Info.read(p)
	if debug {
		fmt.Println("------------------------------------------")
	}
	for i := 0; uint16(i) < h.Info.Count; i++ {
		temp := txdTexture{}
		temp.read(f, p)
		h.Textures = append(h.Textures, temp)
		if debug {
			fmt.Println("------------------------------------------")
		}
	}

	if debug {
		fmt.Println()
	}
	return true
}

func (h *txdInfo) read(p *bufio.Reader) bool  {
	h.Header.read(p)
	buf := make([]byte, h.Header.chunkSize)
	numBytesRead, err := p.Read(buf)
	if uint32(numBytesRead) != h.Header.chunkSize {
		return false
	}
	check(err)
	filePosition += uint32(numBytesRead)

	h.Count = readUint16(buf[:2])
	h.Unknown = readUint16(buf[2:])

	if debug {
		fmt.Printf("	Textures count: %d\n", h.Count)
		fmt.Printf("	Unknown: %d\n", h.Unknown)
	}

	return true
}

func (h *txdTexture) read(f *os.File, p *bufio.Reader) bool  {
	h.Header.read(p)
	h.Data.read(f, p)
	return true
}

func (h *txdTextureData) read(f *os.File, p *bufio.Reader) bool  {
	h.Header.read(p)
	buf := make([]byte, 92)
	numBytesRead, err := p.Read(buf)
	if uint32(numBytesRead) != 92 {
		return false
	}
	check(err)

	h.Version		= readUint32(buf[:4])
	h.FilterFlags	= readUint32(buf[4:8])
	h.TextureName	= string(buf[8:40])
	h.AlphaName		= string(buf[40:72])
	h.AlphaFlags	= readUint32(buf[72:76])
	h.TextureFormat	= string(buf[76:80])
	h.Width			= readUint16(buf[80:82])
	h.Height		= readUint16(buf[82:84])
	h.Depth			= uint8(buf[84])
	h.MipmapCount	= uint8(buf[85])
	h.TexcodeType	= uint8(buf[86])
	h.Flags			= uint8(buf[87])
	h.Palette		= uint8(buf[88])
	h.DataSize		= readUint32(buf[88:92])

	h._DataStart	= filePosition + 92
	h._DataEnd		= filePosition + h.DataSize + 92
	filePosition	+= uint32(h.Header.chunkSize)

	p.Reset(f)
	_, err = f.Seek(int64(filePosition), 0)
	check(err)

	if debug {
		fmt.Printf("	Version:		%d\n", h.Version)
		fmt.Printf("	FilterFlags:	%x\n", h.FilterFlags)
		fmt.Printf("	TextureName:	%s\n", h.TextureName)
		fmt.Printf("	AlphaName:		%s\n", h.AlphaName)
		fmt.Printf("	AlphaFlags:		%x\n", h.AlphaFlags)
		fmt.Printf("	TextureFormat:	%s\n", h.TextureFormat)
		fmt.Printf("	Width:			%d\n", h.Width)
		fmt.Printf("	Height:			%d\n", h.Height)
		fmt.Printf("	Depth:			%d\n", h.Depth)
		fmt.Printf("	MipmapCount:	%d\n", h.MipmapCount)
		fmt.Printf("	TexcodeType:	%d\n", h.TexcodeType)
		fmt.Printf("	Flags:			%d\n", h.Flags)
		fmt.Printf("	Palette:		%d\n", h.Palette)
		fmt.Printf("	DataSize:		%d\n", h.DataSize)
		fmt.Printf("	_DataStart:		%d\n", h._DataStart)
		fmt.Printf("	_DataEnd:		%d\n", h._DataEnd)
	}

	return true
}


/* TODO: 32bit & DXT3 */
func (h *txdFile) replaceTexture(f *os.File, textureId uint16) error {
	condition :=
	 	h.Textures[textureId].Data.Width >= 32 &&
		h.Textures[textureId].Data.Height >= 32 &&
		h.Textures[textureId].Data.Width < 2048 &&
		h.Textures[textureId].Data.Height < 2048

	if condition {
		switch h.Textures[textureId].Data.TextureFormat {
			case "DXT1":
				{
					path := "../src/dds/" +
						strconv.Itoa(int(h.Textures[textureId].Data.Width)) +
						"x" +
						strconv.Itoa(int(h.Textures[textureId].Data.Height)) +
						".dds"

					o, err := os.Open(path)
					if err != nil {
						return err
					}
					bytesf, err := ioutil.ReadAll(o)
					o.Close()

					bytesf = bytesf[0x80:]

					_, err = f.WriteAt(bytesf, int64(h.Textures[textureId].Data._DataStart))
					if err != nil {
						return err
					}
				}
			default:
				return nil
		}
	}

	return nil
}

func (h *txdFile) replaceAll(f *os.File) error {
	var i uint16
	var err error
	for i = 0; i < h.Info.Count; i++ {
		err = h.replaceTexture(f, i)
	}
	return err
}
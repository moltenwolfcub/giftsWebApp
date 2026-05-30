package config

import (
	"encoding/hex"
	"log"
)

const defaultPort = 8040

type Config struct {
	Port int

	//auth - changing any of these will cause all prior hashes to fail
	SaltLength      uint
	ArgonIterations uint32
	ArgonMem        uint32
	ArgonThreads    uint8
	HashLength      uint32
	Pepper          []byte
}

func New() *Config {
	c := &Config{
		Port:            defaultPort,
		SaltLength:      16,
		ArgonIterations: 4,
		ArgonMem:        64 * 1024,
		ArgonThreads:    2,
		HashLength:      32,
	}

	pepper := "fab145e31a6d21c2ee158e7f66195475b8a439d61ca0190bbd449f665ad3da4e"
	//TODO read this in from environment variable rather than hard-coded
	bytePepper, err := hex.DecodeString(pepper)
	if err != nil {
		log.Fatal("Failed to convert pepper")
	}
	c.Pepper = bytePepper

	return c
}

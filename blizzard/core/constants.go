package core

import "github.com/matthewhartstonge/argon2"

var HashConfig = argon2.Config{
	Version:     argon2.Version13,
	Mode:        argon2.ModeArgon2id,
	Parallelism: 1,
	MemoryCost:  1 << 12,
	SaltLength:  1 << 3,
	HashLength:  1 << 5,
	TimeCost:    3,
}

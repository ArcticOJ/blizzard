package main

import (
	"crypto/rand"
	"fmt"
	"github.com/spf13/cobra"
)

var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	root := &cobra.Command{
		Use:   "generator",
		Short: "blizzard key",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, e := cmd.Flags().GetUint8("length")
			if e != nil {
				return e
			}
			b := make([]byte, l)
			if _, e := rand.Read(b); e != nil {
				return e
			}
			for i, j := range b {
				b[i] = charset[j%byte(len(charset))]
			}
			fmt.Println(string(b))
			return nil
		},
	}
	root.Flags().Uint8P("length", "l", 32, "length of string")
	_ = root.Execute()
}

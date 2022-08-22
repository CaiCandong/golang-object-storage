package metadata

import (
	"fmt"
	"testing"
)

func TestMetadata_Get(t *testing.T) {
	var metadata Metadata
	err := metadata.Get("abc.txt", 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(metadata)
}

func TestMetadata_Put(t *testing.T) {
	metadata := &Metadata{
		Name:    "abc.txt",
		Version: 1,
		Size:    200,
		Hash:    "SHA-256=9e0a95c42e3763a0b31a057f3213eeb6",
	}
	err := metadata.Put()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(metadata)
}

// Exists
func TestMetadata_Exists(t *testing.T) {
	metadata := &Metadata{
		Name:    "abc.txt",
		Version: 1,
		Size:    200,
		Hash:    "SHA-256=9e0a95c42e3763a0b31a057f3213eeb6",
	}
	ok := metadata.Exists()
	fmt.Println(ok)
}

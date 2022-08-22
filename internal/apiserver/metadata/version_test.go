package metadata

import (
	"fmt"
	"testing"
)

func TestMetadata_AddVersion(t *testing.T) {
	var metadata Metadata
	err := metadata.AddVersion("abc.txt", 201, "SHA-256=9e0a95c42e3763a0b31a057f3213eeb8")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(metadata)
}

// SearchAllVersions
func TestMetadata_SearchAllVersions(t *testing.T) {
	var metadata Metadata
	versions, err := metadata.SearchAllVersions("abc.txt", 1, 4)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(versions)
}

func TestMetadata_SearchLatestVersion(t *testing.T) {
	var metadata Metadata
	err := metadata.SearchLatestVersion("abc.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(metadata)
}

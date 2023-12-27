package repackage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cbguder/books/epubcheck"
	"github.com/cbguder/books/overdrive"
	"github.com/cbguder/books/repackage"

	"github.com/stretchr/testify/require"
)

func TestEbook(t *testing.T) {
	testBooks(t, "testdata/books")
}

func TestEbookCorpus(t *testing.T) {
	rootDir := os.Getenv("EBOOK_CORPUS_DIR")
	if rootDir == "" {
		t.Skip("EBOOK_CORPUS_DIR not set")
	}

	testBooks(t, rootDir)
}

func testBooks(t *testing.T, rootDir string) {
	t.Helper()

	ents, err := os.ReadDir(rootDir)
	require.NoError(t, err)

	t.Run("books", func(t *testing.T) {
		for _, ent := range ents {
			if !ent.IsDir() {
				continue
			}

			srcDir := filepath.Join(rootDir, ent.Name())

			t.Run(ent.Name(), func(t *testing.T) {
				t.Parallel()
				testSingleBook(t, srcDir)
			})
		}
	})
}

func testSingleBook(t *testing.T, srcDir string) {
	t.Helper()

	dstFile := filepath.Join(t.TempDir(), "book.epub")
	repackageEbook(t, srcDir, dstFile)

	result, err := epubcheck.Check(dstFile)
	require.NoError(t, err)

	c := result.Checker
	if c.NFatal > 0 || c.NError > 0 || c.NWarning > 0 || c.NUsage > 0 {
		for _, msg := range result.Messages {
			t.Logf(
				"%s(%s): %s",
				msg.Severity,
				msg.ID,
				msg.Message,
			)
		}

		t.FailNow()
	}
}

func repackageEbook(t *testing.T, srcDir, dstFile string) {
	t.Helper()

	openbookPath := filepath.Join(srcDir, "_d", "openbook.json")

	openbook, err := overdrive.ReadOpenbook(openbookPath)
	require.NoError(t, err)

	err = repackage.Ebook(srcDir, dstFile, openbook)
	require.NoError(t, err)
}

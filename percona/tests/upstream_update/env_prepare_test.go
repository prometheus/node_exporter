package percona_tests

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestPrepareExporters extracts exporter from client binary's tar.gz
func TestPrepareUpdatedExporter(t *testing.T) {
	if doRun == nil || !*doRun {
		t.Skip("For manual runs only through make")
		return
	}

	if url == nil || *url == "" {
		t.Error("URL not defined")
		return
	}

	prepareExporter(*url, updatedExporterFileName)
}

func extractExporter(gzipStream io.Reader, fileName string) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	exporterFound := false
	for !exporterFound {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if strings.HasSuffix(header.Name, "postgres_exporter") {
				outFile, err := os.Create(fileName)
				if err != nil {
					log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
				}
				defer outFile.Close()
				if _, err := io.Copy(outFile, tarReader); err != nil {
					log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
				}

				exporterFound = true
			}
		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %d in %s",
				header.Typeflag,
				header.Name)
		}
	}
}

func prepareExporter(url, fileName string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	extractExporter(resp.Body, fileName)

	err = exec.Command("chmod", "+x", fileName).Run()
	if err != nil {
		log.Fatal(err)
	}
}

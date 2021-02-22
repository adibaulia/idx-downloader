package main

import (
	"fmt"
	"idx-downloader/http"
	"os"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
)

var (
	oke = "unvr;budi;myor;indf;icbp"
)

func main() {
	year := os.Args[1]
	kodeEmiten := os.Args[2]

	years := strings.Split(year, ",")
	kodeEmitens := strings.ReplaceAll(kodeEmiten, ",", ";")
	fmt.Printf("Downloading Laporan Keuangan dengan kode emiten = %v ditahun = %v\n", strings.ToUpper(kodeEmiten), year)
	var count int
	for _, y := range years {
		reqs := &http.Request{
			Year:       y,
			ReportType: "rdf",
			Periode:    "audit",
			KodeEmiten: kodeEmitens,
		}
		response, err := http.GetLinks(reqs)
		if err != nil {
			panic(err)
		}

		for _, f := range response.Results {
			filename := fmt.Sprintf("FinancialStatement-%v-Tahunan-%v.pdf", reqs.Year, f.KodeEmiten)
			fmt.Printf("Downloading dengan kode emiten %v\n", f.KodeEmiten)
			for _, files := range f.Attachments {
				if files.FileName == filename {
					count++
					downloader("./downloaded/", files.FilePath)
				}
			}
		}
	}
	fmt.Printf("%v files downloaded at ./downloaded\n", count)
}

func downloader(dst, url string) {
	// create client
	_, err := os.Stat(dst)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dst, 0755)
		if errDir != nil {
			fmt.Fprintf(os.Stderr, "Download failed: %v\n", errDir)
			os.Exit(1)
		}

	}
	client := grab.NewClient()
	req, _ := grab.NewRequest(dst, url)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	// Output:
	// Downloading http://www.golang-book.com/public/pdf/gobook.pdf...
	//   200 OK
	//   transferred 42970 / 2893557 bytes (1.49%)
	//   transferred 1207474 / 2893557 bytes (41.73%)
	//   transferred 2758210 / 2893557 bytes (95.32%)
	// Download saved to ./gobook.pdf
}

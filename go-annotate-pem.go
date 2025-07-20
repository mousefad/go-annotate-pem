package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

// RE for finding CERTIFICATE blocks
var certRx = regexp.MustCompile(
	`(?ms)^-----BEGIN CERTIFICATE-----\s*\n.*?\n-----END CERTIFICATE-----\s*$`)

// Argument settings set from command line options
var args struct {
	inPlace bool
	space   string
	paths   []string
}

func main() {
	parseCommandLine()
	for _, path := range args.paths {
		handleFile(path)
	}
}

func handleFile(path string) {
	var outFile *os.File
	var inPath string
	if !args.inPlace {
		outFile = os.Stdout
		inPath = path
	} else {
		backupPath := fmt.Sprintf("%s~", path)
		inPath = backupPath
		log.Printf("moving %q -> %q and re-creating original with annotations",
			path,
			backupPath)
		var originalFileInfo os.FileInfo
		var err error
		if originalFileInfo, err = os.Stat(path); err != nil {
			log.Fatalf("could not get file into for %q", path)
		}
		if err := os.Rename(path, backupPath); err != nil {
			log.Fatalf("failed to backup file %q, aborting", path)
		}
		if outFile, err = os.Create(path); err != nil {
			log.Fatalf(
				"failed to create new file after backup to %q, aborting",
				backupPath)
		}
		defer outFile.Close()
		if err = os.Chmod(path, originalFileInfo.Mode()); err != nil {
			log.Printf("couldn't set the file permissions for: %q", path)
		}
	}

	// Slurp it all up
	buf, err := ioutil.ReadFile(inPath)
	if err != nil {
		log.Fatalf("failed to read %q\n", path)
	}
	matches := certRx.FindAllIndex(buf, -1)
	prevEnd := 0
	for _, match := range matches {
		preText := buf[prevEnd:match[0]]
		certText := buf[match[0]:match[1]]
		fmt.Fprintf(
			outFile,
			"%s%s%s%s",
			preText,
			args.space,
			annotateCert(certText),
			certText)
		prevEnd = match[1]
	}
	if prevEnd < len(buf) {
		fmt.Fprintf(outFile, "%s", buf[prevEnd:])
	}
}

// Take a slice of bytes defining a certificate, and extract information about
// that certificate, returning a textual description as a slice of bytes (one
// field per line, separated with \n). Fields that are extracted: 
// * Subject
// * Issuer
// * Not Before
// * Not After
// * Subj. Alt. Names [if present in cert]
func annotateCert(certText []byte) []byte {
	block, _ := pem.Decode(certText)
	if block == nil {
		return []byte("[unable to parse PEM block]\n")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return []byte(fmt.Sprintf("[unable to parse cert: %v]\n", err))
	}
	var annotations bytes.Buffer
	fmt.Fprintf(&annotations, "Subject:          %s\n", cert.Subject.String())
	fmt.Fprintf(&annotations, "Issuer:           %s\n", cert.Issuer.String())
	fmt.Fprintf(&annotations, "Not Before:       %s\nNot After:        %s\n",
		cert.NotBefore.Format("2006-01-02 15:04:05 -07:00"),
		cert.NotAfter.Format("2006-01-02 15:04:05 -07:00"))
	var sanData bytes.Buffer
	sanCount := 0
	for _, dns := range cert.DNSNames {
		if sanCount > 0 { fmt.Fprintf(&sanData, ",") }
		sanCount += 1
		fmt.Fprintf(&sanData, " DNS:%s", dns)
	}
	for _, email := range cert.EmailAddresses {
		if sanCount > 0 { fmt.Fprintf(&sanData, ",") }
		sanCount += 1
		fmt.Fprintf(&sanData, " Email:%s", email)
	}
	for _, ip := range cert.IPAddresses {
		if sanCount > 0 { fmt.Fprintf(&sanData, ",") }
		sanCount += 1
		fmt.Fprintf(&sanData, " IP:%s", ip.String())
	}
	for _, uri := range cert.URIs {
		if sanCount > 0 { fmt.Fprintf(&sanData, ",") }
		sanCount += 1
		fmt.Fprintf(&sanData, " URI:%s", uri.String())
	}
	if sanData.Len() > 0 {
		fmt.Fprintf(&annotations, "Subj. Alt. Names:%s\n", sanData.String());
	}
	return annotations.Bytes()
}

func parseCommandLine() {
	inPlace := flag.Bool(
		"i",
		false,
		"edit files in place rather than output to stdout")
	space := flag.Bool(
		"s",
		false,
		"add space between certs to help readability")
	os.Args[0] = "annotate_pem"
	flag.Parse()
	args.inPlace = *inPlace
	if *space {
		args.space = "\n"
	}
	args.paths = flag.Args()
}

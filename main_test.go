package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)


func TestRun(t *testing.T) {
  testCases := []struct{
    name string
    root string
    cfg config
    expected string
  }{
    { name:"NoFilter", root: "testdata", cfg: config{ ext: "", size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
    { name:"FilterExtensionMatch", root: "testdata", cfg: config{ ext: ".log", size: 0, list: true}, expected: "testdata/dir.log\n"},
    { name:"FilterExtensionNoMatch", root: "testdata", cfg: config{ ext: ".gz", size: 0, list: true}, expected: ""},
    { name:"FilterExtensionSizeMatch", root: "testdata", cfg: config{ ext: ".log", size: 10, list: true}, expected: "testdata/dir.log\n"},
    { name:"FilterExtensionSizeNoMatch", root: "testdata", cfg: config{ ext: ".log", size: 20, list: true}, expected: ""},
  }

  for _, tc := range testCases {
    t.Run(tc.name, func(tt *testing.T) {
      var buffer bytes.Buffer

      if err := run(tc.root, &buffer, tc.cfg); err != nil {
        t.Fatal(err)
      }

      res := buffer.String()

      if tc.expected != res {
        t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
      }
    })
  }
}

func TestRunDelExtension(t *testing.T) {
  // Create the test cases we going to run
  testCases := []struct {
    name string
    cfg config
    extNoDelete string
    nDelete int
    nNoDelete int
    expected string
  }{
    { name: "DeleteExtensionNoMatch", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 0, nNoDelete: 10, expected: ""},
    { name: "DeleteExtionMatch", cfg: config{ext: ".log", del: true}, extNoDelete: "", nDelete: 10, nNoDelete: 0, expected: ""},
    { name: "DeleteExtensionMixed", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 5, nNoDelete: 5, expected: ""},
  }

  // Loop through all test cases defined above
  for _, tc := range testCases {
    // Call run for each case
    t.Run(tc.name, func(tt *testing.T) {
      var (
        buffer bytes.Buffer
        logBuffer bytes.Buffer
      )

      tc.cfg.wLog = &logBuffer

      // Create files and temp dir for this case
      tempDir, cleanup := createTempDir(tt, map[string]int{
        tc.cfg.ext: tc.nDelete,
        tc.extNoDelete: tc.nNoDelete,
      })
      // call function to remove all files
      // once the test is done
      defer cleanup()

      if err := run(tempDir, &buffer, tc.cfg); err != nil {
        tt.Fatal(err)
      }

      res := buffer.String()

      if tc.expected != res {
        tt.Errorf("Expected %q, got %q instead\n", tc.expected, res)
      }

      filesLeft, err := os.ReadDir(tempDir)

      if err != nil {
        tt.Error(err)
      }

      if len(filesLeft) != tc.nNoDelete {
        tt.Errorf("expected %d files left, got %d instead \n", tc.nNoDelete, len(filesLeft))
      }

      expLogLines := tc.nDelete + 1
      // split the bytes by new line
      lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
      if len(lines) != expLogLines {
        t.Errorf("Expected %d log lines, got %d instead\n", expLogLines, len(lines))
      }
    })
  }
}

func TestRunArchive(t *testing.T) {
  // Define the test cases to run
  testCases := []struct {
    name string
    cfg config
    extNoArchive string
    nArchive int
    nNoArchive int
  }{
    { name: "ArchiveExtensionNoMatch", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 0, nNoArchive: 0},
    { name: "ArchiveExtensionMatch", cfg: config{ext: ".log"}, extNoArchive: "", nArchive: 10, nNoArchive: 0},
    { name: "ArchiveExtensionMixed", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 5, nNoArchive: 5},
  }

  // Loop through the test cases
  for _, tc := range testCases {
    t.Run(tc.name, func (tt *testing.T) {
      var buffer bytes.Buffer
      // Create a temp dir to store the files
      tempDir, cleanup := createTempDir(tt, map[string]int{
        tc.cfg.ext: tc.nArchive,
        tc.extNoArchive: tc.nNoArchive,
      })
      // cleanup after done
      defer cleanup()

      // create new dir to 
      // copy the archived files to
      archiveDir, cleanupArchive := createTempDir(tt, nil)
      defer cleanupArchive()

      // assign it to the config
      tc.cfg.archive = archiveDir

      // Run the program
      if err := run(tempDir, &buffer, tc.cfg); err != nil {
        tt.Fatal(err)
      }

      // Look for the files that were actually archived
      pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tc.cfg.ext))
      expFiles, err := filepath.Glob(pattern)
      if err != nil {
        tt.Fatal(err)
      }

      expOut := strings.Join(expFiles, "\n")
      res := strings.TrimSpace(buffer.String())

      // check if expected output
      // is the actual output that got generated
      if expOut != res {
        tt.Errorf("Expected %q, got %q instead \n", expOut, res)
      }

      filesArchived, err := os.ReadDir(archiveDir)
      if err != nil {
        tt.Fatal(err)
      }

      // check that the archived file count
      // is the amount we expected
      if len(filesArchived) != tc.nArchive {
        tt.Errorf("Expected %q, got %q instead \n", tc.nArchive, len(filesArchived))
      }
    })
  }
}

// Create temp dir for testing
// instead of hardcoding the dir
func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
  // makes it a helper function
  t.Helper()

  // Create directory
  tempDir, err := os.MkdirTemp("", "walktest")
  if err != nil  {
    t.Fatal(err)
  }

  // loop through each key
  for k, n := range files {
    // Create as many files as passed in
    for j := 1; j <= n; j++ {
      // create file name using key and current index
      fname := fmt.Sprintf("file%d%s", j, k)
      // path of the dir to write to
      fpath := filepath.Join(tempDir, fname)
      // create the file
      if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
        t.Fatal(err)
      }
    }
  }

  return tempDir, func() { os.RemoveAll(tempDir) }
}

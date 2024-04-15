package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)


type config struct {
  // extension to filter out
  ext string
  // min file size
  size int64
  //list files
  list bool
  // delete file
  del bool
  // Logger
  wLog io.Writer
  //archive directory
  archive string
}

func main() {
  root := flag.String("root", ".", "Root directory to start")
  logFile := flag.String("log", "", "Log deletes to this file")
  // Action options
  list := flag.Bool("list", false, "List files only")
  del := flag.Bool("del", false, "Delete file")
  archive := flag.String("archive","", "Archive directory")
  // Filter options
  ext := flag.String("ext", "", "File extension to filter out")
  size := flag.Int64("size", 0, "Minimum file size")
  flag.Parse()


  var (
    f = os.Stdout
    err error
  )

  if *logFile != "" {
    f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

    if err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }

    defer f.Close()
  }

  c := config{
    ext: *ext,
    size: *size,
    list: *list,
    del: *del,
    archive: *archive,
    wLog: f,
  }

  if err := run(*root, os.Stdout, c); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func run(root string, out io.Writer, cfg config) error {
  delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)
  return filepath.Walk(root, func (path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    if filterOut(path, cfg.ext, cfg.size, info) {
      return nil
    }

    if cfg.list {
      return listFile(path, out)
    }

    if cfg.archive != "" {
      if err := archiveFile(cfg.archive, root, path); err != nil {
        return err
      }
    }

    if cfg.del {
      return delFile(path, delLogger)
    }


    return listFile(path, out)
  })
}

package main

import (
  "fmt"
  "os"
  "net/http"
  "time"
  "strconv"
)

var pid int

type DownloadFile struct {
  Filename string
  Filesize int64
  Offset int64
  File *os.File
  DoSeek bool
  LastSize int64
} 

func (f *DownloadFile) procRunning(pid int) (bool) {
  proc, err := os.FindProcess(pid)
  if err != nil || proc == nil {
    return false
  } else {
    return true
  }
}

func (f *DownloadFile) haveData(size int) (bool) {
  var ret = false
  s, _ := f.File.Stat()
  currSize := s.Size()
  if currSize == f.LastSize {
    if pid > 0 && f.procRunning(pid) {
      ret = false
    } else {
      ret = true
    }
  } else if (f.Offset + int64(size)) < currSize {
    ret = true
  } else {
    f.LastSize = currSize
  }
  return ret
}

func (f *DownloadFile) Read(p []byte) (n int, err error) {
  size := len(p)
  for !f.haveData(size) {
    time.Sleep(2000 * time.Millisecond)
  }
  if f.DoSeek {
    f.Offset, err = f.File.Seek(f.Offset, os.SEEK_SET)
    if err != nil {
      fmt.Println(err)
    }
    f.DoSeek = false
  }
  n, err = f.File.Read(p)
  f.Offset += int64(n)
  return n, err
}

func (f *DownloadFile) Seek(offset int64, whence int) (int64, error) {
  if whence == os.SEEK_SET {
    f.Offset = offset
  } else if whence == os.SEEK_CUR {
    f.Offset += offset
  } else {
    f.Offset = f.Filesize + offset
  }
  f.DoSeek = true
  return f.Offset, nil
}

func NewDownloadFile(filename string) (*DownloadFile) {
  f := &DownloadFile{}
  f.Filename = filename
  f.Filesize = 1000000000000
  var err error
  f.File, err = os.Open(filename)
  if err != nil {
    fmt.Println(err)
  }
  if err != nil {
    fmt.Println(err)
  }
  return f
}

var fname string

func handler(w http.ResponseWriter, r *http.Request) {
  file := NewDownloadFile(fname) 
  defer file.File.Close()
  http.ServeContent(w, r, fname, time.Time{}, file)
  
}

func main() {
  numargs := len(os.Args)
  if !(numargs == 2 || numargs == 3) {
    fmt.Println("Usage: dlsrv filename [pid]")
    os.Exit(1)
  }
  fname = os.Args[1]
  if numargs == 3 {
    pid, _ = strconv.Atoi(os.Args[2])
  }

  for i:=0;i<20;i++  {
    if _, err := os.Stat(fname); os.IsNotExist(err) {
      time.Sleep(1000 * time.Millisecond)
    } else {
      break
    }
  }
  time.Sleep(4000 * time.Millisecond)
  fmt.Println("Listening on: 127.0.0.1:9696") 
  http.HandleFunc("/", handler)
  http.ListenAndServe(":9696", nil)
}

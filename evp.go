package main
 
// #include <stdlib.h>
// #include <openssl/rsa.h>
// #include <openssl/engine.h>
// #include <openssl/evp.h>
// #include <openssl/pem.h>
// #include <openssl/bio.h>
// #cgo CFLAGS: -I/usr/local/opt/openssl/include
// #cgo LDFLAGS: /usr/local/Cellar/openssl/1.0.1e/lib/libcrypto.a
// #//////cgo LDFLAGS: -L/usr/local/opt/openssl/lib -lcrypto
import "C"
import (
    "os"
    "flag"
    "fmt"
    "log"
    "unsafe"
    "io/ioutil"
)

type M struct {
    //ctx *_Ctype_RSA
    rsa_key *C.RSA
}

func Init() {
    C.OPENSSL_add_all_algorithms_noconf()
}

func Cleanup() {
    C.EVP_cleanup()
}

func LoadRSAKey(file_path string, passwd string) *M {
  m := &M{}
 //   m.ctx = C.RSA_new()
  bufPEM, err := ioutil.ReadFile(file_path)
  if err != nil {
      log.Fatal(err)
  }
  //m.rsa_key = C.RSA_generate_key(2048, 3, nil, nil)
  //log.Print(m.rsa_key)
  //buf := C.create_string_buffer(raw, len(raw))
  //buf := C.CString("hihi")
  //buf := []byte(password)
  dataLen := len(bufPEM)
  bio := C.BIO_new_mem_buf((unsafe.Pointer(&bufPEM[0])), (C.int(dataLen)))

  buf := make([]byte, 200)
  //log.Print(C.BIO_read(bio, (unsafe.Pointer(&buf[0])), C.int(dataLen)))
  //C.BIO_ctrl(bio, 1, 0, nil)  // rewind
  //log.Print(string(buf))
  //log.Print("Pw: ", string(passwd), len(passwd))
  cs := C.CString(passwd)
//  m.rsa_key, err = C.PEM_read_bio_RSAPrivateKey(bio, nil, nil, (unsafe.Pointer(&passwd)))
  m.rsa_key, err = C.PEM_read_bio_RSAPrivateKey(bio, nil, nil, (unsafe.Pointer(cs)))
  C.BIO_free(bio)
  errNo := C.ERR_get_error()
  if errNo != 0 {
    //C.BIO_reset(bio)
    log.Print("Load rsa key error: ", errNo)
    C.ERR_error_string_n(errNo, (*C.char)(unsafe.Pointer(&buf[0])), 200)
    log.Fatal(string(buf))
  }
  log.Print(errNo)
  if m.rsa_key == nil {
    log.Fatal("rsa key error: ", errNo)
  }

  log.Print(m.rsa_key)
  return m
}

func main() {
  flag.Parse()
  args := flag.Args()
  if len(args) < 1 {
    fmt.Println("Input file is missing.");
    os.Exit(1);
  }
  if len(args) < 2 {
    fmt.Println("Password is missing.");
    os.Exit(1);
  }
  Init()
  m := LoadRSAKey(args[0], args[1])
  print (m)
  Cleanup()
}

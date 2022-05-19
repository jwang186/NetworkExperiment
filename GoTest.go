   package main

   import (
       "fmt"
       "os"
   )

   type AdminModeType string

   const (
       UDS AdminModeType = "uds"
       TCP AdminModeType = "tcp"
   )

   func setEnv() {
    os.Setenv("TEST", "hh")
    defer os.Unsetenv("TEST")
   }

   func main() {
       os.Setenv("TEST", "one")
       setEnv()
       fmt.Println(os.Getenv("TEST"))
   }
package main

import (
  ".."
  "fmt"
  "flag"
)

func main() {
  secret := flag.String("secret", "", "The shared secret")
  url := flag.String("url", "", "The URL to encode")
  flag.Parse()

  digest := camo.NewDigest(*secret)
  fmt.Println(digest.Calculate(*url))
}
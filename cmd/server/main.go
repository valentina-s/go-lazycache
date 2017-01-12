package main

import "fmt"
import "net/http"
import "net/url"
import "strings"
import "sort"

import "github.com/amarburg/go-lazycache"

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/"

func main() {

  url,err := url.Parse( OOIRawDataRootURL )
  fs, err := lazycache.OpenHttpFS( *url )

  if err != nil {
    panic( fmt.Sprintf("Error opening HTTP FS Source: %s", err.Error() ) )
  }

  serverAddr := "localhost:5000"


  //http.HandleFunc("*.mov/*", lazycache.MoovHandler )

  // Reverse hostname
  splitHN := strings.Split( fs.Uri.Host, "." )
  fmt.Println(splitHN)

  for i, j := 0, len(splitHN)-1; i < j; i, j = i+1, j-1 {
      sort.StringSlice(splitHN).Swap(i,j)
  }

  root := fmt.Sprintf("/%s%s", strings.Join(splitHN,"/"), fs.Uri.Path )
fmt.Println(root)
  http.Handle(root, lazycache.MakeTreeHandler( fs, root ) )
  http.HandleFunc("/", lazycache.Index )

  fmt.Printf("Starting http handler at http://%s/\n", serverAddr)
  fmt.Printf("Fs at http://%s%s\n", serverAddr, root )

  http.ListenAndServe(serverAddr, nil)
}

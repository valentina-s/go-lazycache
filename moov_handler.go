package lazycache

import "fmt"
import "net/http"
import "strings"
import "strconv"

//import "strings"
import "encoding/json"
import "image/png"

import "github.com/amarburg/go-lazyfs"
import "github.com/amarburg/go-lazycache/quicktime_store"
import "github.com/amarburg/go-lazyquicktime"

func MoovHandler(node *Node, path []string, w http.ResponseWriter, req *http.Request) {
	//  fmt.Fprintf( w, "Quicktime handler: %s with residual path (%d): (%s)\n", node.Path, len(path), strings.Join(path,":") )

	lqt, have := quicktime_store.HaveEntry(node.trimPath)

	if !have {
		uri := node.Fs.Uri
		uri.Path += node.Path
		fmt.Println(uri.String())
		fs, err := lazyfs.OpenHttpSource(uri)
		if err != nil {
			http.Error(w, "Something's went boom opening the HTTP Soruce!", 500)
			return
		}

		lqt, err = quicktime_store.AddEntry(node.trimPath, fs)
		if err != nil {
			http.Error(w, fmt.Sprintf("Something's went boom parsing the quicktime file: %s", err.Error()), 500)
			return
		}
	}

	if len(path) == 0 {
		b, err := json.MarshalIndent(lqt, "", "  ")
		if err != nil {
			fmt.Fprintln(w, "JSON error:", err)
		}

		w.Write(b)
	} else {
		switch strings.ToLower(path[0]) {
		case "frame":
			handleFrame(node, lqt, path[1:], w, req)
		default:
			http.Error(w, fmt.Sprintf("Didn't understand request \"%s\"", path[0]), 500)
		}
	}

}

func handleFrame(node *Node, lqt *lazyquicktime.LazyQuicktime, path []string, w http.ResponseWriter, req *http.Request) {

	if len(path) == 0 {
		http.Error(w, fmt.Sprintf("Need to specify frame number"), 500)
    return

	}

	frameNum, err := strconv.Atoi(path[0])
  if err != nil {
    http.Error(w, fmt.Sprintf("Error parsing frame number \"%s\"", path[0]), 500)
    return
  }

	img, err := lqt.ExtractFrame(frameNum)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating image for frame %d: %s", frameNum, err.Error()), 500)
    return
	}

	err = png.Encode(w, img)
}

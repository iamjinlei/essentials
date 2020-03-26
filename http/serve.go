package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/phayes/freeport"

	"github.com/iamjinlei/memfs"
)

func ServeImage(img []byte, w, h int) error {
	fsMap := map[string][]byte{
		"/snap.png": img,
		"/index.html": []byte(fmt.Sprintf(`
<!doctype html>
<html>
	<head>
		<title>Selenium debug snapshot</title>
		<link rel="icon" href="data:;base64,iVBORw0KGgo=">
	</head>
	<body>
		<img src="snap.png" style="width:%vpx;height:%vpx" alt="snap.png">
	</body>
</html>
`, w, h)),
	}

	var wg sync.WaitGroup
	wg.Add(3)

	fs, err := memfs.New(fsMap, map[string]func(path string){
		"Close": func(path string) {
			fmt.Printf("close %v\n", path)
			wg.Done()
		},
	})
	if err != nil {
		return err
	}

	port, err := freeport.GetFreePort()
	if err != nil {
		return err
	}

	fmt.Printf("serving http://localhost:%v\n", port)

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(fs)))
	srv := &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: mux}

	go func() {
		srv.ListenAndServe()
	}()

	wg.Wait()
	return srv.Shutdown(context.TODO())
}

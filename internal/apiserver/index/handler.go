package index

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		`<html>
	    <head>
	        <meta charset="utf-8"></meta>
	        <title>Uploader</title>
	    </head>
	    <body>
	        <form action="/objects" method="post" enctype="multipart/form-data">
	            <input type="file" id="file" name="file">
	            <input type="submit" name="submit" value="upload">
	        </form>
	    </body>
	</html>`)
}

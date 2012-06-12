// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html"
)

func htmlTemplate(requestedHash string, formatedHashes string) string {
	return fmt.Sprintf(`<!DOCTYPE HTML>
<html>
<body>

<h1>%s</h1>
<ul>
%s
</ul>
</body>
</html>`, html.EscapeString(requestedHash), formatedHashes)
}

func htmlFormatHash(hash string, password string) string {
	return fmt.Sprintf("<li><code>%s</code> <b>%s</b></li>\n", html.EscapeString(hash), html.EscapeString(password))
}

package common

const (
	COMMAND_START = 1

	NOT_FOUND = `
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html><head>
<title>404 Not Found</title>
</head><body>
<h1>Not Found</h1>
<p>The requested URL #?# was not found on this server.</p>
</body></html>
`

	METHOD_NOT_ALLOWED = `
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html><head>
<title>405 Method Not Allowed</title>
</head><body>
<h1>Method Not Allowed</h1>
<p>The #?# method is not allowed for the requested URL.</p>
</body></html>
`
	FORBIDDEN = `
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html><head>
<title>403 Forbidden</title>
</head><body>
<h1>Access forbidden</h1>
<p>You don't have permission to access the requested resource.</p>
</body></html>
`
)

var (
	Port        = 80
	ContextPath = ""
	BasicAuth   = ""
	LocalAddr   = ""
	RemoteAddr  = ""
	Command     int
)

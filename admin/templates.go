package main

import (
	"html/template"
)

const indexTpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>rex admin</title>
<style type="text/css">
* {
	padding: 0px;
	margin: 0px;
	box-sizing: border-box;
}

body {
	font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
	font-size: 24px;
	font-style: normal;
	font-variant: normal;
	font-weight: 500;
	line-height: 26.3999996185303px;

	background: #eee;
	color: #222;
}

.wrap {
	max-width: 600px;
	margin: 10px;
}

input, button {
	display: block;
	float: left;
	height: 30px;
}

input {
	width: 90%;
	padding: 5px;
}

button {
	width: 10%;
}

.clear {
	clear: both;
}

ul {
	margin-top: 10px;
	border-top: 1px solid #bbb;
}

li a {
	display: block;
	width: 100%;
	padding: 5px;
	font-size: 16px;
	color: #222;
	text-decoration: none;
	border-bottom: 1px solid #ccc;
}

li a:hover {
	color: #333;
	background: #ddd;
}
</style>
</head>
<body>
<div class="wrap">
<form method="post" action="/">
<input type="text" placeholder="new script" name="name"><button type="submit">Create</button>
</form>
<div class="clear"></div>
<ul>
{{range .}}
<li><a href="/{{.ID}}">{{.ID}}</a></li>
{{end}}
</ul>
</div>
</body>
</html>`

var index = template.Must(template.New("index").Parse(indexTpl))

const editTpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>{{.ID}} - rex admin</title>
<style type="text/css">
* {
	padding: 0px;
	margin: 0px;
	box-sizing: border-box;
}

body {
	font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
	font-size: 24px;
	font-style: normal;
	font-variant: normal;
	font-weight: 500;
	line-height: 26.3999996185303px;

	background: #eee;
	color: #222;
}

.clear {
	clear: both;
}

.wrap {
	max-width: 600px;
	margin: 10px;
}

input, button, textarea {
	display: block;
}

input, textarea {
	width: 100%;
}

button {
	width: 25%;
	padding: 5px;
	font-size: 16px;
	text-align: center;
	border: 1px solid #ccc;
	border-radius: 3px;
	color: #222;
	background: transparent;
	margin-top: 15px;
}

input, button {
	height: 38px;
	padding: 5px;
}

input {
	margin-bottom: 15px;
}

textarea {
	min-height: 400px;
	padding: 5px;
}

a {
	display: block;
	width: 25%;
	float: right;
	padding: 5px;
	font-size: 16px;
	color: #222;
	text-decoration: none;
	border: 1px solid #ccc;
	border-radius: 3px;
	text-align: center;
	margin-bottom: 15px;
}

a:hover, button:hover {
	color: #333;
	background: #ddd;
	cursor: pointer;
}

</style>
</head>
<body>
<div class="wrap">
<a href="/">Go back</a>
<div class="clear"></div>
<form method="post" action="/{{.ID}}">
<input type="text" name="name" value="{{.ID}}">
<textarea name="code">{{.Code}}</textarea>
<button type="submit">Update</button>
</form>
</div>
</body>
</html>`

var edit = template.Must(template.New("edit").Parse(editTpl))

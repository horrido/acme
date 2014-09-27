<!DOCTYPE html>

<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8" />
<title>StarDust by Free Css Templates</title>
<meta name="keywords" content="" />
<meta name="description" content="" />
<link href="/static/css/default.css" rel="stylesheet" type="text/css" />
</head>

<body>
<!-- start header -->
<div id="header-bg">
	<div id="header">
		<div align="right">{{if .InSession}}
		Welcome, {{.First}} [<a href="http://localhost:8080/user/logout">Logout</a>|<a href="http://localhost:8080/user/profile">Profile</a>]
		{{else}}
		[<a href="http://localhost:8080/user/login/home">Login</a>]
		{{end}}
		</div>
		<div id="logo">
			<h1><a href="#">StarDust<sup></sup></a></h1>
			<h2>Designed by FreeCSSTemplates</h2>
		</div>
		<div id="menu">
			<ul>
				<li class="active"><a href="http://localhost:8080/home">home</a></li>
				<li><a href="#">photos</a></li>
				<li><a href="#">about</a></li>
				<li><a href="#">links</a></li>
				<li><a href="#">contact </a></li>
			</ul>
		</div>
	</div>
</div>
<!-- end header -->
<!-- start page -->
<div id="page">
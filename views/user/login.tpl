<div id="content">
<h1>Please login</h1>
&nbsp;
{{if .flash.error}}
<h3>{{.flash.error}}</h3>
&nbsp;
{{end}}
{{if .Errors}}
{{range $rec := .Errors}}
<h3>{{$rec}}</h3>
{{end}}
&nbsp;
{{end}}
<form method="POST">
<table>
<tr>
    <td>Email address:</td>
    <td><input name="email" type="text" autofocus /></td>
</tr>
<tr>      
    <td>Password:</td>
    <td><input name="password" type="password" /></td>
</tr>
<tr><td>&nbsp;</td></tr>
<tr>
    <td>&nbsp;</td><td><input type="submit" value="Login" /></td><td><a href="http://localhost:8080/user/register">Register</a></td>
</tr>
</table>
</form>
</div>

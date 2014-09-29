<div id="content">
<h1>Please register</h1>
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
    <td>First name:</td>
    <td><input name="first" type="text" autofocus /></td>
</tr>
<tr>
    <td>Last name:</td>
    <td><input name="last" type="text" /></td>
</tr>
<tr>
    <td>Email address:</td>
    <td><input name="email" type="text" /></td>
</tr>
<tr>      
    <td>Password (must be at least 6 characters):</td>
    <td><input name="password" type="password" /></td>
</tr>
<tr>      
    <td>Confirm password:</td>
    <td><input name="password2" type="password" /></td>
</tr>
<tr><td>&nbsp;</td></tr>
<tr>
    <td>&nbsp;</td><td><input type="submit" value="Register" /></td>
</tr>
</table>
</form>
</div>

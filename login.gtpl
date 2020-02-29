<html>
<head>
<title></title>
</head>
<body>
<form method="post" action="/login">
	用户名:<input type="text" name="username">
	密码:<input type="password" name="password">
	age: <input type="text" name="age">
	<select name="fruit">
		<option value="apple">apple</option>
		<option value="pear">pear</option>
		<option value="banana">banana</option>
	</select>
	<input type="hidden" name="token" value="{{.}}">
	<input type="submit" value="登陆">
</form>
</body>
</html>

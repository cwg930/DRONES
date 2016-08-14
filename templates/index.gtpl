<html>
  <head>
    <title></title>
  </head>
<body>
  <form action="/users" method="post">
    Name:<input type="text" name="username">
    Age:<input type="number" name="age">
    <input type="submit" value="submit">
  </form>
  <br>
  <form enctype="multipart/form-data" action="/files" method="post">
    File:<input type="file" name="uploadFile">
    <input type="submit" value="Upload">
  </form>
</body>
</html>

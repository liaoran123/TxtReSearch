
ajaxpost=function(url,JSON,callback){ 
  var xhr=new XMLHttpRequest();
  xhr.onreadystatechange=function(){
      if(xhr.readyState==4){
          if(xhr.status>=200&&xhr.status<=300||xhr.status==304){   
           // var objs = JSON.stringify(xhr.responseText);        
              callback(xhr.responseText)
          }
      }
  }
    // 拼接JSON数据，比如我们的参数{"id":1,"name":"小明","age":18}
   // 转换为id=10001&name=小明&age=18
   var temp = [];
   for(var k in JSON) {
     temp.push(k+"="+encodeURI(JSON[k]));
   }
   var str=temp.join("&")||null
   xhr.open("post",url,true);
   xhr.setRequestHeader('Content-Type','application/x-www-form-urlencoded');
   //xhr.setRequestHeader('Content-Type','application/json');
   xhr.send(str)
 
}
//--------cookie----------------
const cookies = {
  cookie: {
    // 设置cookie
    set: (name, value, day) => {
      const date = new Date();
      date.setDate(date.getDate() + day);
      document.cookie = name + "=" + value + ";expires=" + date;
    },
    // 获取cookie
    get: (key) => {
      var arr = document.cookie.split("; ");
      for (var i = 0; i < arr.length; i++) {
        var arr1 = arr[i].split("=");
        if (arr1[0] == key) {
          return arr1[1]
        }
      }
      return ""
    },
    // 删除cookie
    remove: (name) => {
      cookie.set(name, '', -1)
    }
  }

}
/*
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>封装cookie</title>
</head>
<body>
  
</body>
<script src="./cookie.js"></script>
<script>
  const {cookie}= cookies; //引入声明一下
  cookie.set('maomin', '22', '0.5');
  console.log(cookie.get("maomin"));
</script>
</html>

*/
//-----字符串截取----------
var getmindstr = function r(con, l, r, ll, rl) {
  //function getmindstr(con, l, r, ll, rl) {
  //--获取字符中间的字符,ll,rl是否最后一个匹配
  var lp, rp, cp;
  lp =
    l == ""
      ? 0
      : ll == false
        ? (lp = con.indexOf(l))
        : (lp = con.lastIndexOf(l));
  if (lp == -1) return "";
  lp = lp + l.length;
  rp =
    r == ""
      ? con.length
      : rl == false
        ? con.indexOf(r, lp)
        : (rp = con.lastIndexOf(r));
  if (rp == -1) return "";
  cp = rp - lp;
  if (cp < 0) return "";
  return con.substr(lp, cp);
}
var getCaptcha = function (btnobj, imgobj) {
  const timestamp = Date.parse(new Date());
  document.getElementById(btnobj).setAttribute("data-id", timestamp)
  document.getElementById(imgobj).setAttribute("src", "/Captcha/?id=" + timestamp)
}  
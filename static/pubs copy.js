//公共函数
var sev = "http://127.0.0.1:8049";//"http://research.soufoshuo.com";//
//var jjsev = "http://jianjie.soufoshuo.com";// "http://127.0.0.1:8050";
var l = "";
var wh = window.location.href;
if (wh.indexOf("l=") != -1) {
  l = getmindstr(wh + "&", "l=", "&", true, false);
}
if (l != "") l = "&l=1";
//----选项卡触发事件--------
$(function () {
  //$("#spinner").hide(); //隐藏加载动画
  //------搜索按钮--------
  $("#param").keydown(function (e) {
    if (e.keyCode == 13) {
      $("#ssbt").click();
    }
  });
  $("#ssbt").click(function () {
    var kw = $("#param").val().trim();
    if (kw == "") return;
    var t = $("#param").attr("data-param");
    var pl = l;
    if (pl == "") pl = "0";
    var jlv = $("#bttext").text();
    var dir = $("#bttext").attr("data-jlv");
    var url = "/s/" + kw + "?dir=" + dir + l;
    window.location.href = url;
  });
  setvalue = function (ci) {
    $("#param").val(ci);
    scrollTo(0, 0);
  };
});
//搜索提示框
$(function () {
  //打开提示框
  $("#param").click(function () {
    $("#OpenexampleModal").click()
    //input是隐藏状态，或者未在页面内创建input元素节点，需要等待节点正常显示后使用focus才能生效
    setTimeout(function () {
      $('#searchinput').focus();
    }, 490)
  })
  //动态加载搜索相关词
  $("#searchinput").on("input", function () {    
    var myReg = /^[\u4e00-\u9fa5]+$/; //过滤非中文
    let inval = $("#searchinput").val()    
    if (myReg.test(inval)) {   
      $("#leftlist").attr("data-p","")   
      $("#righlist").attr("data-p","")     
      getIdxpfx("leftlist","")
      getIdxpfx("righlist","")
    }
  })
});
var getIdxpfx = function(id,add){  
  let inval = $("#searchinput").val()
  let caid = $("#bttext").attr("data-jlv");
  let obj=$("#"+id)
  let isr=obj.attr("data-isr")
  let p=$("#"+id).attr("data-p")  
  let url = sev + "/api/Idxpfx/?kw=" + inval + "&caid=" + caid+"&isr="+isr+"&p="+p ;
  $.getJSON(url, function (data) {
    console.log(data)
    let htmls = "";
    for (let index = 0; index < data.kwpxf.length; index++) {
      let element = data.kwpxf[index];
      element=element.replace(new RegExp(inval , "g"), "<b>"+inval+"</b>")
      htmls += "<li class='list-group-item'>" + element + "</li>"; //<li class="list-group-item">A fourth item</li>
    }     
    if (add==""){ 
      obj.html(htmls)       
    }       
    else{//“加载更多” 
      let ohtml=obj.html()
      obj.html(ohtml+htmls)  
    }      
    obj.attr("data-p",data.p)     
  })
}
$(function () {
  //选项卡
  var cjlv = $.cookie("jlv");
  if (typeof cjlv != "undefined") {
    var jlvs = cjlv.split("|");
    $("#bttext").text(jlvs[0]);
    $("#bttext").attr("data-jlv", jlvs[1]);
  }
  $("#jlvDropdown .dropdown-item").click(function () {
    var jlvtext = $(this).text();
    var jlv = $(this).attr("data-jlv");
    $("#bttext").text(jlvtext);
    $("#bttext").attr("data-jlv", jlv);
    $.cookie("jlv", jlvtext + "|" + jlv, {
      expires: 7,
      path: "/",
    });
  });
});


$(function () {
  $("#param").click(function () {
    $("#OpenexampleModal").click() //打开搜索提示框模块
    //input是隐藏状态，或者未在页面内创建input元素节点，需要等待节点正常显示后使用focus才能生效
    setTimeout(function () {
      $('#searchinput').focus();
    }, 490)

  })
  //--返回顶部
  $("#retop").click(function () {
    scrollTo(0, 0);
  });

  $("#clearcooke").click(function () {
    $.removeCookie("sousuolishi");
    $("#sousuolishi-Htmls").html("");
  });
});
//------播放文字转语音---------

$(function () {
  $("#play").click(function () {
    window.speechSynthesis.resume();
    var text = $("#jingwennr").text();
    text = text.replace(/\s+/g, "，");
    text = text.replace(/\s+/g, "，");
    var utterThis = new window.SpeechSynthesisUtterance(text);
    window.speechSynthesis.speak(utterThis);
  });

  $("#pause").click(function () {
    window.speechSynthesis.pause();
  });
});

//---cookie----
var rdcookie = function (cname, cont) {
  var fsc = $.cookie(cname);
  if (typeof fsc == "undefined") fsc = "";
  var ofcs = fsc.replace(new RegExp(cont + ",", "g"), ""); //   fsc.replace(cont + ",", "")
  $.cookie(cname, cont + "," + ofcs, {
    expires: 365,
    path: "/",
  });
  if (ofcs == $.cookie(cname)) {
    //--cooke已经到达最大长度。
    ofcs = getmindstr("." + ofcs, ".", "|", false, true);
    ofcs = getmindstr("." + ofcs, ".", "|", false, true); //--删除最后两个
    $.cookie(cname, cont + "," + ofcs, {
      expires: 365,
      path: "/",
    });
  }
};
/*
var getlishicooks = function () {
  var fsc = $.cookie("sousuolishi");
  if (fsc == null)
    fsc =
      "阿弥陀佛,极乐世界,般若,禅定,灭尽三昧,五戒,云何,何者是,贪嗔痴,戒定慧,贪嗔 痴,戒定 慧,《";
  var arr = fsc.split(",");
  var fs = "",
    url,
    title;
  for (var i = 0; i < arr.length; i++) {
    if (arr[i] == "") continue;
    title = decodeURI(arr[i]);
    fs = fs + '<option value="' + title + '">';
  }
  return fs;
}

//设置cooke
var scooke = function (name, value) {
  $.cookie(name, value, {
    expires: 21,
    path: '/'
  });
}
var delcooke = function (name) {
  $.cookie(name, {
    expires: -1,
    path: '/'
  });
}
*/
//获取验证码
/*
var getCaptcha=function (){
  btnobj=$("#check-btn")
  imgobj=$("#capimg")
  getCaptchaToObj(btnobj,imgobj)
}*/
//将验证码参数赋予对应的2个组件
var getCaptchaToObj = function (btnobj, imgobj) {
  const timestamp = Date.parse(new Date());
  $(btnobj).attr("data-id", timestamp)
  $(imgobj).attr("src", "/Captcha/?id=" + timestamp)
}
//计算时间差
var getDifferenceTime = function (sj) {
  let nowTime = Date.now();//获取当前时间对应的毫秒数                
  let sjstr = sj.replace("/-/g", "/") //字符串转时间
  let eightTime = new Date(sjstr)//getDate($("fbsj").text())
  let differenceTime = nowTime - eightTime;
  let day = Math.floor(differenceTime / 1000 / 24 / 3600)
  return day
}
//-----字符串截取----------
var getmindstr = function (con, l, r, ll, rl) {
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
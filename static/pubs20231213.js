
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
    //var jlv = $("#bttext").text();    
    var dir = $("#bttext").attr("data-dir");
    var url = "/s/" + kw + "_" + dir;
    if (dir == "10000") url = "/taobao/?kw=" + kw + "&p=1";
    window.location.href = url;
  });
  setvalue = function (ci) {
    $("#param").val(ci);
    scrollTo(0, 0);
  };
});

$(function () {
  //下拉框绑定cook数据
  $("#datas").append(getlishicooks());
  //-------绑定下拉框数据-----------
  $("#param").on("input", function () {
    var inval = $("#param").val();
    if (inval == "") {
      //空值时，加载cook数据
      $("#datas").children().filter("option").remove();
      $("#datas").append(getlishicooks());
      return;
    }
    var myReg = /^[\u4e00-\u9fa5]+$/; //过滤非中文
    //console.log(inval.substr(0,1))
    if (myReg.test(inval) || inval.substr(0, 1) == "《") {
      //var name = $("#param").val();
      $("#datas").children().filter("option").remove();
      var jlv = $("#bttext").text();
      var caid = $("#bttext").attr("data-dir");
      debugger
      if (caid == "10000") return
      var url = "/idxpfx/?kw=" + inval + "&caid=" + caid;//"/getidxkey/?kw=" + inval + "&caid=" + caid + l;
      $.getJSON(url, function (data) {
        //console.log(data)
        var htmls = "";
        for (let index = 0; index < data.length; index++) {
          const element = data[index];
          if (element.key.indexOf(" ") != -1) {
            //过滤有空格的项
            continue;
          }
          htmls += "<option>" + element.key + "</option>";
        }
        if (htmls != "") {
          $("#datas").children().filter("option").remove();
          $("#datas").append(htmls);
        }
      });
    }
  }); /*.focus(function () {
   //   $("#datas").children().filter("option").remove();
   // });*/
});

$(function () {
  //选项卡  
  if (window.location.href.indexOf("taobao") != -1) {
    $("#bttext").attr("data-dir", "10000");
    $("#bttext").text("淘宝");
    return
  }
  if (window.location.href.indexOf("/dir/") != -1) {
    return
  }
  var cjlv = $.cookie("jlv");
  if (typeof cjlv != "undefined") {
    var jlvs = cjlv.split("|");
    $("#bttext").text(jlvs[0]);
    $("#bttext").attr("data-dir", jlvs[1]);
  }

  $("#jlvDropdown .dropdown-item").click(function () {
  var jlvtext = $(this).text();
  var jlv = $(this).attr("data-dir");
  $("#bttext").text(jlvtext);
  $("#bttext").attr("data-dir", jlv);
  $.cookie("jlv", jlvtext + "|" + jlv, {
    expires: 21,
    path: "/",
  });
});
});


$(function () {
  //--返回顶部
  $("#retop").click(function () {
    scrollTo(0, 0);
  });
  //清除cooke
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

var getlishicooks = function () {
  var fsc = $.cookie("sousuolishi");
  if (fsc == null)
    fsc =
      "论语,道德,般若,黄鹤楼,卦象,经脉,《";
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
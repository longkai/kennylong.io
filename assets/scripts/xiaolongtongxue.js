function ajax(method, url, success, fail) {
  var ajax = window.XMLHttpRequest ? new XMLHttpRequest() : new ActiveXObject("Microsoft.XMLHTTP");
  ajax.onreadystatechange = function () {
    if (ajax.readyState == XMLHttpRequest.DONE) {
      if (ajax.status == 200) {
        success(JSON.parse(ajax.responseText));
      } else {
        fail(ajax.status);
      }
    }
  };
  ajax.open(method, url, true);
  ajax.send();
}

function daysAgo(s) {
  return Math.floor((new Date().getTime() - new Date(s).getTime()) / (60 * 1000 * 60) / 24);
}

Array.prototype.forEach.call(document.querySelectorAll('.social-btn'), function (el) {
  openLink(el, el.getAttribute(('href')));
});

function openLink(el, url) {
  el.addEventListener('click', function () {
    window.open(url, '_self');
  });
}

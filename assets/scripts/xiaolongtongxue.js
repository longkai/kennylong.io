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

function sinceNow(s) {
  return Math.floor((new Date().getTime() - new Date(s).getTime()) / (60 * 1000 * 60) / 24);
}

Array.prototype.forEach.call(document.querySelectorAll('.social-btn'), function (el) {
  var platform = el.getElementsByTagName('span')[0].textContent.toLowerCase();
  console.log(platform);
  switch (platform) {
    case 'gmail': // :)
      openLink(el, 'mailto:' + atob('aW0ubG9uZ2thaUBnbWFpbC5jb20='));
      break;
    case 'stackoverflow':
      openLink(el, 'https://stackoverflow.com/users/3280791/longkai');
      break;
    case 'github':
      openLink(el, 'https://github.com/longkai');
      break;
    case 'instagram':
      openLink(el, 'https://www.instagram.com/xiangchiyoutiao');
      break;
    case 'medium':
      openLink(el, 'https://medium.com/@longkai');
      break;
    case 'twitter':
      openLink(el, 'https://twitter.com/xiaolongtongxue');
      break;
  }
});

function openLink(el, url) {
  el.addEventListener('click', function () {
    window.open(url, '_self');
  });
}
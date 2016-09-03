var icons = JSON.parse('{"settings":{"layout":"compact","prefix":"sprite","padding":5,"uri":true,"stylesheet":"css"},"canvas":{"sprites":[{"name":"github","src":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAADFUlEQVRoQ+3ZS6iVVRQH8N/tBZalhApSIoSNMrRABEFU7G1Da2QRGTkT8RGOBBFE8AHqSKmgEqEHNMkswahBRpPSHjNrUhg9ENNeiJYs2AfOvXz3nn2/fXfnCHdNz9pr/f/fWmfv9RhyncvQdY7fJIF+R3C0CMzBXjyG2/sM8hI+wCb8MBJLE4G7cQZ39hn4SPfnsQA/dv/QROAtPDVg4Dtw3sbTvQj8jjsGlMBFTOtF4L8BBd+BNSxrmlJokkDlCE5GID7wAXyClViNWS2/+i94ByexDOsz7ExIBObj2+RsCjZiK6biKr7Hb7iQdKZjBu7BjfgDu7APfyed+/DN/0HgL9zW4Gg25uHLBLAJSxB8AGfxU4PCn7i1B4niCMQ7EV+0hkTEht3zDU6KCfybUqUT+okiEqkYqXVD7QiE/YWpXpoo8GEn6pzTGQaLI3AY6zIctVE5hBdrRiBCHDfJr23QZZyZmW6w+LOPJkUROIJnMoCUqLyBNbUIvIBXStBlnF2Ll2sReDx1Rxk4WqtEF3i8FoFHcaI1tLyDj+DDWgSeReRoTQkfr9UisAV7aqLHZuyuReAYnqxM4D2sqkUgiq0onaOgqyFRyMUbM1ZBV/QOBOht2FEDfbK9vYftYgLxGi/B1xNM4n6cSoXiWKaLCYTxqOWfyCy+cng+mO7+nM6uFYGo0yM37+1CczndFgfxcw7KBp27UhsZHd1NmTZaEYhr7SUswtHUeXX8BZGP8DG+wqeIAVSTRLOyApEuy7EUN2cC76i1IhCHX02lbrSOXyAqx5ESBBbjn1FA3YLPECnTVloT6L6B4rl/PzXo3UByir2Ybb7ZFj3DdxrjnczF/R/9QOT8Q9iZmvT4j3yO59I0Yix8EcFz/SIQfvdjQwGAGKtcKThflELhN+Y+z+P1AhAl89diAh3c8R+IDu279EXn4t1MUgNBoAlr7tKwKoGSBUdtAlkLjpIVU20CWSumtku+0WamTakWujGJG49kL/nCaJs1a/Sx0ZDnSPTVD+coYtxr1ky7g6GWm7ODgbYBxSSBfofmGmrenzFHYaoqAAAAAElFTkSuQmCC"},{"name":"gmail","src":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAB5ElEQVRoQ+2ZrU4EMRSFvw1I/oKGB0AhcHgSHIEEXgGDwyKxOAyvAAkExxPgEDg8aAg/EgI5m2kyadrtzLRldjet2Wy69+ecc+/tbGfAhK/BhOdPAdC3gkWBokAkA1NbQmvAKbAFzEeSFGv+CdwBx8Cz7cylgJK/B5ZiIye2fwM2gae6XxeAa2A3cfBU7m6AvRCAjzEoGx9g5bYYAvBrWRuVVoBzYCcVnR4/t8AR8FLt+/IZbrtKaKRBBUBABCjlUsJKXADqKzkAOZ+rppQCzkSi+KmUPQG+HL6yADBxNoALQJ9d1gNwCOjTt7ICUFApICV0bkiZJktMi3GVohQYtbIDMMGbNrndpCHAyQGoXEZJrinlanJfk9YBuHwnB/AdaDq7yfVdgHxNav9+NvcUMow0ZVT5tFHMHu3JFbAdtq3pUM/8OwAl1GaqhKZWLwAMq6G53uTc6BWAgLhO1jYnd+8AjBqmyc0UavrsNDYAQgeUb78ACD1O2/tdmS4KiIEuf2iKAoGaK01cmjhyLJUSegcWIlnMZd7oYusS2M+VQaTfK+Cg7sN1Duih6xFYjgyW2vwVWK/d2A39+94PrAJnwPYY3JO2vl5PzVxWf1P7hiYraymdFwVSstnFV1GgC2spbSZegT8TrYoxdOJQhQAAAABJRU5ErkJggg=="},{"name":"instagram","src":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAACo0lEQVRoQ+3aS6hNYRQH8N9FBuSRGVFekSjGGFJCimJkYMDAo7zmBjL1Ko8BAwMjihIyMMSYIo+8ipgRYiCvVp1d++72ufY+Z3/3nFtnDc/5vrX+/7W+b621195DxrgMjXH8BgR6HcF2EViCY1iLKT0G+Q13cBjviljKCAT4B5jeY+BF85+xEs/yf5QRuIbNfQY+g3MdW/5H4GsfHJt2/gts0/5H4G+PvL8a9wu2V+Fe4bdhp6bsCPWKQLuEUsQzIJDqhA0ikHm2V3dgzF/iqkezby9x3xB4ggu4i7ctVHOxBruwtCrSNuuSReAnDuE8/rQxPg67cQITOySShECA39DyehVcEY1bHZJIQmAvzlVBnluzB2dr7onljRN4jBX4XRPMeDzEspr7GidwEKdqgsiWH8DJmntrEyiW+JdYkDMaHozM04lERooIZvIKCwuKum7migS+Y1LOyFTEY18nEo+r0eNn8gOTUxMoPvA0SSAcEfry0ngEnmNRoiP0AotTE4j8vT5npMlLfLtVT5JGIHL+mZyFJtPovpLa0PgRmoPXmJAj0UQh+4X5JbOfxgkE7mjWduYINNFKXGw1e8VsloTALDwtZIxumrnIbDFQ+1CSi5MQCDubEIOm6DDzUredjs41Bmk32hSSZATC3v5Wa1wkUbWoBfhowU+PsCEpgSwSlzuY5kXR2j6C5zNOyQmEoZk4ih2F7FTm2Mg2l3AEHyuEalQIZDhmt+7GRsxDpNyQGIu/wc2Wx99XAD6qEaiBp/bSUY1AbXQVNgwIlE3MKjiukSUdjde/lPTkjaBpQEmlFxxXsLUBYylUXMW2vOKykXakwkeYkQJBFzo/YTmGpeB2M/nI38exroMK2wXG0q21X7M2DSCpvsGnBkndW0H5mI/APzJaozE53lmmAAAAAElFTkSuQmCC"},{"name":"medium","src":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAADt0lEQVRoQ+2YaciOWRjHf3bJyJIU+cBQg6xZsiSELNM0amrsywdC2T5I2YksSaFkKWEmpqaU0pSGkki2kH1JvliiiGSZsfWv6+h+b/dz7/f7eus59Xx4nuec61y/cy3nuk4NqvmoUc31pwxQ1RYsW6BsgYwnEORCW4D6wCbgbkb5pZZL/t9AG2CXfd6l2SsI4CNQE/gAHADWFACyEljhUfgRsAHYDbxNAhIE8NknIG+Q9sBVoF6AogLZaBaJBRIHwO2TF8i/wNCIU44NkgQgD5CxwMEELvLYXEtxEmiRNABpQRoBN4GWCQDc1CfA+iDXygLgBVlivhum21ZgTgrlvUv2AVO9P+QBIHn/A3VDlOsBnANqZQB4A3QCHhQBIJmlyhKl5DNA7wzKa+liYJ1fRl4WCAOYBWzPqLxipxugO0qfr6MIgIbAa9uhBXALaJwBQPfSEOAEMNt/GEUA6EZdZArvBaZkUF5L/wAmA+4wmhRtgf+A7sB1A1H6SzteAD8BT4E/gQn+WCvCAjL5SWAQUBu4DHRMSaD42QEMBo6b8hV0LgpA+sp19gMDzX+TNk9Ku32BOnYIssQ3yaJIgGdm/udA0lhQpukFXLL0ubZU4ikSQHvuBGYCzS0bNY3pSrq151m/cA1oUFUAn4B+wFlghgFFMagS7QC8Ao4Ao30LKi0G3L5yA7mDgvuU+XUYxDjgL2AMcChgYqUDSIe5wDa7Tc9bdgqCUK8wHNBleANo/b0AyB2URVTfbwYWBCimnriLta/qyhaWMFOVWEC6yC3kHj9YX9DKp+AqQL1yZ+Cipc8ghsIB7gDqe4PGMOAY8Ju9Srg590zx93YJDggJksIBdPuOB5YB7XyKCE5uIkX/AUba/yOAo8A0YE9EmiocwO0vENUuS30gAtNTzY+AcvxhQL1yM3Mt3Rlho9IAvCATAbWdsoiCVZ3VfWC+udJD63enR10ScYo597AVQ1aFKVG1jiwiEFnktu+CUs1zOqSr824UaQFdNv2Tah9zc4l1IFes1tH3C0DXmHtGAqhxUPGlwEoyoixQSlbcEsOtjwRwE/sAy4FRMSnSAqhcVselGNFjb9SIDeAEqY4RyM9JskOUFgH/C2SSxUgYSGIAt1dPD0jQaae1gJ9FMeEs0jYANDWAk6VHKlnkF1/g5gXgTb/OIg5E/XaFV+0sm+qdRiC/2nN53CyS1Ltc1vrdHshWl8ypSSXbfBVlLz1vQSnFpFuWxQLpdsx5VRkg5wNNLK5sgcRHlvOCam+BL1qzwjEJeGhwAAAAAElFTkSuQmCC"},{"name":"stackoverflow","src":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAADu0lEQVRoQ+2YWchNURiGn988ZHZjJkO5UOaZyIVZCBFluHBhKJKhZMiQDMWFEkq5MStj5MIQmYcolEiGDClj5gy9WquWbe+zz9777PM7db76+0/7rLW/913f8H7rlFDgVlLg+CkSiBnB2sDrmHv/2pbvCNQEtgHNgG7Al6Qk8kmgI7DHgBfuTcC0QiLQBrgIVHFAjwH2JiGRzwgI51RgswP4HdAOeBiXRL4JCOcOYJwD+CrQA/gWh0SaBHoBAvfZA6y6ed7Seb4BmP0/EegCnATuAsOAZx5wSpsLQEXz/BcwAjgYlUQaEWhowOm/TOBF4poH3HRgo/NMutAeeBSFRBoEBgH7gMoOEKXRRJ+Osx8Y6aw7D/QBvmdLIg0C8t0JOADUd4AoTZYCywF9lknYbgBNnXWrgQWlTUD+lUIi0cEDZjcw2SnuzsBZoIJTD4OBY9mQSCsC1rdEazswygPmCjDcKe45wDpnzSugrU/x/8MpbQJyKB9KnUXmswXx1JBQcWvNIWCIg/AM0A/4kSkSuSBQBjgNHDfzzZsAh2PNIOcW90dgkin6OqYeGjn7VxjigRxyQUDpYeeZD8AWQMKkE/ZaUHEvBlYaRT4FlDOFvgpYmHYEpLbeQtVYoJFhLXDHAyCouHcCU4wiS5UnACfCCjlpBNQ5dHoai2v5OFO7PAyoNarHhxX3JaPIP4GXYeBtgWWzLmxNNTNpzjLt02/9OWANcMSkh19xq/t0B+6HObTfJ42A148iMh6YC7QOAHHbpJZSTIpri7ss0N80hGzxp3apV2caCsw3V0c/QE+A9cBWQ7Y5sCtr5GZhriPg57+3ITLQowN2rYa4maboo+JPFAGNBJ+Ay4CU9WbIEKYr5TyTMmqTrumCr+tmZIsbgUrAe6C84/GrESJLSKTuOYObXdoE0OigllnV5HzfyMgTplBXM/OH+X3rREiERO652VQXmAFoZJB4xbK4EVC3kNq2AlSwUUzFa8loCnX1Icp7/qzNRGC087YHwHWft9cA9HuPRmL9aVRoEAGF5icVt6we0NPZK92w0Qp8ZSYC9tKhzVF+hBIBEbGkRFBE/WwZsMR8oTuARM6aJtOjYYeRBgGvT6WYUs1GSP8160v0pBUW9H9LwO8QBV4k9KuFJlhZXglIMVuEhTeL75UmL0qDwC1AwpTUpAmPiwQCjjFTFyr4FEqaOn7781rERQI+J5BqBHQxl7SnaY2BAWkpcZrAg96d01GiSCDGCSSOQAyf+d8S90KTf6QBHosESjsUvwEPr9Qxtjy0KwAAAABJRU5ErkJggg=="},{"name":"twitter","src":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAC10lEQVRoQ+3aS+jNWxQH8I/3xCukrsdESCaICQaU66aES7lidqNEysCMTO41MMDMVUrJjesxuDchN4+IMlFihIFHHikhJHlctOp38nfueWznv4/z/+vsOqPf2mt9v2utvfdae58euvno0c3xaxNodQSrRWA0tmMeBrQY5EucxAbcK8dSicAoXMWQFgMvN/8MM3C944dKBA5jaRcDX4LzN5bUI/AcA7sogRcYVI/Axy4KvgTri6yplEJtAk2OYDsCTXZwXfUtjcBr/IuzeIM4c5ZjbAfY/+FVjZ0wC4HTmIxhdf31WSDmrMLdsjl9sQNzcAL7cQaDq+jOQmAm4ogPUMMTSFzAj3hbR7YX/sJC9GsmgaF4igk4inE1gAXoMXiQQHQuRha/Lc0kEEYeFgbi1I7CbyUVy/NT+CkBfElkMQ4iUqvSyJJCC3CsTPtUbER869Ph226sTiSwAvsQqVRtZCFwBL9UsRALez6mYyKuYH0igUM19GYtJWKb+7lYxInYksTOYVYdySwRiEUZCzi2u3o7SxLyQugGxn8LArElBvicI8rkx2Xrp2mLOJqeZTnRYxH+SdCZJYXCzjrsTDCYKhK61iYIZyMQtvbid9xKMFxLpD9uJ5Ym2QhE0RX79iNcxIdOkIgbh22J87MRCHtx3G9ONFxNLLx/Ez8k6slKIKLwK/5MNF5J7EBRUqeqyEqgZDSuO/bgUlHkpYJZgz9ShQu5phAI3deKgu5yIqDoDXahd6J8SSwrgbinOV94/zjeJ4Dpid+wKUG2KQdZNDFbcafYPr/mGmZ20X1NaRB8TMsSgSeIMjl+QaTWiOYnCr84uaNh6ezIQqAEIrwfl62xeGMrjCr1HUYUXVg069MayPNaJLMS6Kw3G5nfJtCI13LOaUcgpzcb0VU3At3+gaMrPzH97zakOz3yxU3gJNzvmHff5TNrIwurZXPafzVomesrdTetBtOI/U/sB4wx/D9yWwAAAABJRU5ErkJggg=="}]}}').canvas.sprites;

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
  var platform = el.getElementsByTagName('span')[0].textContent.toLowerCase();
  for (var i = 0; i < icons.length; i++) {
    var icon = icons[i];
    if (icon.name === platform) {
      el.style.backgroundImage = 'url("' + icon.src + '")';
      break;
    }
  }
  openLink(el, el.getAttribute(('href')));
});

function openLink(el, url) {
  el.addEventListener('click', function () {
    window.open(url, '_self');
  });
}

function linkfy(cdn, href) {
  if (cdn) {
    return qualifyURL(href);
  }
  return href;
}

// http://stackoverflow.com/questions/470832/getting-an-absolute-url-from-a-relative-one-ie6-issue
function qualifyURL(url) {
  var a = document.createElement('a');
  a.href = url;
  return a.cloneNode(false).href;
}

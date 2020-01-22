const disPassword = 'Please offer your password to access: ' +
    '    <label>\n' +
    '        <input id="password" name="password" type="password" placeholder="password" required/>\n' +
    '    </label>\n' +
    '    <input type="submit" value="Go" onclick="getPassword()">\n';

const header = '<strong><a href="'+urlShort+'">PasteBin </a></strong>' + urlShort +
    '    <label>\n' +
    '        <input id="id" name="id" type="text"  placeholder="Paste\'s number" required/>\n' +
    '    </label>\n' +
    '    <input type="submit" value="Go" onclick="getId()">';

function setHeader() {
    document.getElementById("header").innerHTML = header;
}

let password = null;

function getDisContent(lan, cont) {
    return '<pre class="line-numbers"><code class="language-' + lan + '">' + cont + '</code></pre>';
}

function getContent(url) {
    let httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', url, true);
    console.log("getContent url: " + url);
    httpRequest.send();
    httpRequest.onreadystatechange = function () {
        if (httpRequest.readyState === 4 && httpRequest.status === 200) {
            console.log("getContent: " + httpRequest.responseText);
            let paste = JSON.parse(httpRequest.responseText);
            if (paste.status === "3") {
                console.log("111");
                document.getElementById("box").innerHTML = disPassword;
            }else if (paste.status === "0") {
                console.log("222");
                document.getElementById("box").innerHTML = getDisContent(paste.highlight, paste.content);
                loadScript("js/prism.js", (function() { }));
                //document.getElementById("code").innerText = paste.content;
            }
        }
    };
}

function getId() {
    location.replace(urlShort+document.getElementById("id").value);
}

function getPassword() {
    password = document.getElementById("password").value;
    getContent(url+"&pwd="+password);
}

function getUrlRelativePath() {
    let url = document.location.toString();
    let arrUrl = url.split("//");

    let start = arrUrl[1].indexOf("/");
    let relUrl = arrUrl[1].substring(start);//stop省略，截取从start开始到结尾的所有字符

    if(relUrl.indexOf("?") !== -1){
        relUrl = relUrl.split("?")[0];
    }

    console.log("getUrlRelativePath: " + relUrl);
    return relUrl;
}

function loadScript(url, callback){
    let script = document.createElement ("script");
    script.type = "text/javascript";
    if (script.readyState){ //IE
        script.onreadystatechange = function(){
            if (script.readyState === "loaded" || script.readyState === "complete"){
                script.onreadystatechange = null;
                callback();
            }
        };
    } else { //Others
        script.onload = function(){
            callback();
        };
    }
    script.src = url;
    document.getElementsByTagName("head")[0].appendChild(script);
}





function getCookie(name) {
	var dc = document.cookie;
	var prefix = name + "=";
	var begin = dc.indexOf("; " + prefix);
	if (begin == -1) {
		begin = dc.indexOf(prefix);
		if (begin != 0) return null;
	}
	else                         
	{
		begin += 2;
		var end = document.cookie.indexOf(";", begin);
		if (end == -1) {
			end = dc.length;
		}
	}
	// because unescape has been deprecated, replaced with decodeURI
	// //return unescape(dc.substring(begin + prefix.length, end));
	return decodeURI(dc.substring(begin + prefix.length, end));
}

var myCookie = getCookie("access_token");
if (myCookie == null) {
	document.getElementById("suqi_sypeeoui").textContent = "Login";
}else{
	document.getElementById("suqi_sypeeoui").textContent = "Area Utente";
}                                                                                                                     

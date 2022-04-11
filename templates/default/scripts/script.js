$(document).ready(function(){
	$("#toggle-menu").click(function(){
		var menuStatus = $(document.body).attr("data-menu")
		$(document.body).attr("data-menu", menuStatus == "open"? "close": "open")
	})
})
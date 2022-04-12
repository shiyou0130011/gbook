$(document).ready(function(){
	$("#toggle-menu").click(function(){
		var menuStatus = $(document.body).attr("data-menu")
		$(document.body).attr("data-menu", menuStatus == "open"? "close": "open")
	})

	$("#format-change").click(function(){
		$("#format-menu").toggle()
	})

	var fontSizeList = ["small", "normal", "large"]
	$("#format-menu .text-decrease").click(function(){
		var fontIndex = fontSizeList.indexOf(document.body.dataset.size)
		document.body.dataset.size = fontSizeList[Math.max(fontIndex - 1, 0)]
	})
	
	$("#format-menu .text-increase").click(function(){
		var fontIndex = fontSizeList.indexOf(document.body.dataset.size)
		document.body.dataset.size = fontSizeList[Math.min(fontIndex + 1, fontSizeList.length)]
	})
})
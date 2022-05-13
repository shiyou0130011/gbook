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
		document.body.dataset.size = fontSizeList[Math.min(fontIndex + 1, fontSizeList.length - 1)]
	})
	$("#font-change-sans-serif").click(function(){
		document.body.dataset.font = "sans-serif"
	})
	$("#font-change-serif").click(function(){
		document.body.dataset.font = "serif"
	})
})
$(document).ready(function(){
	$("pre code").each(function(){
		var lang = null
		this.classList.forEach(function(c){
			if (c.indexOf("language-") >= 0){
				lang = c.replace("language-", "")
			}
		})
		if(PR){
			var newHTML = PR.prettyPrintOne(this.innerText, lang, true)
			console.log(newHTML)
			this.innerHTML = newHTML
		}
	})
})
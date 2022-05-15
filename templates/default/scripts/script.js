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
		
		localStorage.fontSize = document.body.dataset.size
	})
	
	$("#format-menu .text-increase").click(function(){
		var fontIndex = fontSizeList.indexOf(document.body.dataset.size)
		document.body.dataset.size = fontSizeList[Math.min(fontIndex + 1, fontSizeList.length - 1)]
		
		localStorage.fontSize = document.body.dataset.size
	})
	$("#font-change-sans-serif").click(function(){
		document.body.dataset.font = "sans-serif"

		localStorage.font = "sans-serif"
	})
	$("#font-change-serif").click(function(){
		document.body.dataset.font = "serif"

		localStorage.font = "serif"
	})
	$(".change-theme").click(function(){
		var theme = this.dataset.targetTheme
		$("link.theme-change-element").each(function(){
			$(this).attr("href", $(this).attr(theme))
		})
		
		localStorage.theme = theme
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
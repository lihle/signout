var comboplete = new Awesomplete('input.dropdown-input', {
	minChars: 0,
});

Awesomplete.$('.dropdown-input').addEventListener("click", function() {
	if (comboplete.ul.childNodes.length === 0) {
		comboplete.minChars = 0;
		comboplete.evaluate();
	} else if (comboplete.ul.hasAttribute('hidden')) {
		comboplete.open();
	} else {
		comboplete.close();
	}
});

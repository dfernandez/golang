$(document).ready(function(){
	// Cookie legal
	if (!Cookies.get("colernio-cookie-legal")) {
		div  = '<nav class="navbar navbar-default navbar-fixed-bottom cookies-legal-navbar" role="navigation">'
		div += '<div class="alert alert-info alert-dismissible cookies-legal col-xs-10 col-xs-offset-1 col-sm-8 col-sm-offset-2" role="alert">'
		div += '<button type="button" class="close close-cookie" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>'
		div += cookieLegalMessage + '</div></div>'

		$('body').append(div);

		// Grab close button to save cookie
	    $('.close-cookie').click(function( e ){
	        Cookies.set('colernio-cookie-legal', true, { expires: 7, path: '/' });
	    });
	}
})
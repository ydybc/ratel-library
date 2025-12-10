(function($) {
	var cse,bds;
	$(function(){
		$('#searchBtn').on('click',function(){
			var x = $('#searchAll').val();
			if(!x){
				$(location).attr('href', '/search');
			}else{
				$(location).attr('href', '/search?q='+$('#searchAll').val());
			}
			
		});
	});
})(jQuery);


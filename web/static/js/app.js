(function($) {
  'use strict';

  $(function() {
    var $fullText = $('.admin-fullText');
    $('#admin-fullscreen').on('click', function() {
      $.AMUI.fullscreen.toggle();
    });

    $(document).on($.AMUI.fullscreen.raw.fullscreenchange, function() {
      $fullText.text($.AMUI.fullscreen.isFullscreen ? '退出全屏' : '开启全屏');
    });
    
  	if($.cookie){
			var setRead = ['fontSize','fontColor','backgroundColor','fontFamily'];
				for(var itme in setRead){
					var setcookie = $.cookie('z_'+setRead[itme]);
					if(setcookie){
						if(setRead[itme]=='backgroundColor'){
							$('body').addClass(setcookie);
						}
						$('#am-read-centent').addClass(setcookie);
					}
				}
  	}


  	$('#am-font-family').change(function(){
  		var setFamily = $(this).val();
  		$('#am-font-family option').each(function(){
  			classRm($(this).val());
  		});
  		classAdd(setFamily);
  		$.cookie('z_fontFamily', setFamily , { path: '/', expires: 10 });  
  	});
  	
  	$('#am-font-color').change(function() {
  		var setColor = $(this).val();
  		$('#am-font-color option').each(function(){
  			classRm($(this).val());
  		});
  		classAdd(setColor);
  		$.cookie('z_fontColor', setColor , { path: '/', expires: 10 });  
  	});
  	
  	$('#am-font-size').change(function() {
  		var setSize = $(this).val();
  		$('#am-font-size option').each(function(){
  			classRm($(this).val());
  		});
  		classAdd(setSize);
  		$.cookie('z_fontSize', setSize , { path: '/', expires: 10 }); 
  	});
  	
  	$('#am-background-color').change(function() {
  		var setBackgroundColor = $(this).val();
  		$('#am-background-color option').each(function(){
  			classRm($(this).val());
  			$('body').removeClass($(this).val());
  		});
  		classAdd(setBackgroundColor);
  		$('body').addClass(setBackgroundColor);
  		$.cookie('z_backgroundColor', setBackgroundColor , { path: '/', expires: 10 }); 
  	});
  	
  	$('.btn-loading-example').click(function () {
		  var $btn = $(this)
		  $btn.button('loading');
		    setTimeout(function(){
		      $btn.button('reset');
		  }, 5000);
		});

		$('#addmark').click(function(){
			var addid = new Object();
			addid.articleid = bookid;
			addid.chapterid = chapterid;
			var $addMarkBtn = $(this);
			$.ajax({
				type:"post",
				url:addUrl,
				dataType:"json",
				timeout:"10000",
				contentType : "application/json",
				data: JSON.stringify(addid),
				beforeSend:function(){
					$addMarkBtn.button('loading');
				},
				error:function(xhr){
					$addMarkBtn.button('reset').text('请重试!');
				},
				success:function(data){
					$addMarkBtn.button('reset').text(data.Info);
					return false;
				}
			});
		});
		if(typeof bookid!="undefined"){
			$('footer').append('<div class="am-modal am-modal-prompt" tabindex="-1" id="collect"><div class="am-modal-dialog">'
				+'<div class="am-modal-hd">提交一点详细内容吧！</div><div class="am-modal-bd"><textarea rows="3" class="am-modal-prompt-input" ></textarea></div><div class="am-modal-footer">'
				+'<span class="am-modal-btn" data-am-modal-cancel>取消</span><span class="am-modal-btn" data-am-modal-confirm>提交</span></div></div></div>');
		}
	$('#collectbutton').on('click', function() {
		var mess = new Object();
		mess.title = $('h2').html();
		mess.articleid=bookid;
		mess.chapterid=chapterid;
		$('#collect').modal({
		  relatedTarget: this,
		  onConfirm: function(e) {
					mess.content=e.data;
					$.ajax({
						type:"post",
						url:"/collect",
						dataType:"json",
						timeout:"5000",
						data:JSON.stringify(mess),
						beforeSend:function(){
							
						},
						error:function(xhr){
							alert('请重试!');
						},
						success:function(data){
							alert(data.Info);	
							return false;
						}
					});
		  },
		  onCancel: function(e) {

		  }
		});
	});
  });
	function classRm(rmVal){
		$('#am-read-centent').removeClass(rmVal);
	}
	function classAdd(addVal){
		$('#am-read-centent').addClass(addVal);
	}

})(jQuery);
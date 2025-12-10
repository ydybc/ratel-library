(function($) {
  'use strict';
  //全屏缩放
  /*
  $(function() {
    var $fullText = $('.admin-fullText');
    $('#admin-fullscreen').on('click', function() {
      $.AMUI.fullscreen.toggle();
    });

    $(document).on($.AMUI.fullscreen.raw.fullscreenchange, function() {
      $fullText.text($.AMUI.fullscreen.isFullscreen ? '退出全屏' : '开启全屏');
    });
  });*/
  $(function() {
	  $('#doc-vld-msg').validator({
	    onValid: function(validity) {
	      $(validity.field).closest('.am-form-group').find('.am-alert').hide();
	    },

	    onInValid: function(validity) {
	      var $field = $(validity.field);
	      var $group = $field.closest('.am-form-group');
	      var $alert = $group.find('.am-alert');
	      // 使用自定义的提示信息 或 插件内置的提示信息
	      var msg = $field.data('validationMessage') || this.getValidationMessage(validity);

	      if (!$alert.length) {
	        $alert = $('<div class="am-alert am-alert-danger"></div>').hide().
	          appendTo($group);
	      }

	      $alert.html(msg).show();
	    }
	  });
	});
  $(function(){
	  if($.cookie){
		  if($.cookie('msize')){
			  nsize($.cookie('msize'));
			  
		  }
		  if($.cookie('mblackboard')){
			  nblackboard($.cookie('mblackboard'));
		  }
	  }
	  //alert($.cookie('msize'));
	  var $novelsize = $('#readtxt');
	  $novelsize.on('change', function() {
		  nsize($(this).val())
	  });
	  var $novelblackboard = $('#readbox');
	  $novelblackboard.on('change', function() {
		  nblackboard($(this).val());

	  });
	  $('#bookmark,#bookmark2').click(function () {
		  var $btn = $(this);
		  var addid = new Object();
			addid.articleid = articleid;
			addid.chapterid = chapterid;
		  $.ajax({
			 type:"post",
			 url:"/addmark",
			 contentType:"application/json",
			 dataType:"json",
			 timeout:"5000",
			 cache:false,
			 data:JSON.stringify(addid),
			 beforeSend:function(){
					$btn.button('loading');
				 },
			 error:function(XMLHttpRequest){
					$btn.button('reset');
					$btn.text('请重试!');
			 },
			 success:function(data){
					$btn.button('reset');
					$btn.text(data.Info);
				 return false;
			 }
		  });
		  //setTimeout($('#bookmark').popover('close'),2000);
	  });
	  $('#delbutton').click(function () {
		  var $btn = $(this);
		  var delid = new Object();
			delid.bookid = new Array();
			$("input[name='checkbox']:checked").each(function() {
				delid.bookid.push(Number($(this).val()));
			});
			console.log(delid.bookid.length)
			if (!delid.bookid.length){
				alert("请选择书籍");
				return false;
			}
		  $.ajax({
			 type:"post",
			 url:"/delmark",
			 contentType:"application/json",
			 dataType:"json",
			 timeout:"5000",
			 cache:false,
			 data:JSON.stringify(delid),
			 beforeSend:function(){
				$btn.button('loading');
				 },
			 error:function(XMLHttpRequest){
				$btn.button('reset');
				$btn.text(data.Info);
			 },
			 success:function(data){
				$btn.button('reset');
				alert(data.Info);
				if(data.Status){
					window.location.reload();
				}
			 }
		  });
		  //setTimeout($('#bookmark').popover('close'),2000);
	  });
	  $('#novel-list').on('change', function() {
		    novel_list_ajax($(this).val());
		  });
	  $("#uppage").on('click',function(){
		  var n = $('#novel-list');
		  var m = n.find("option:selected").prev().val();
		  if(m){
			  n.val(m);
			  novel_list_ajax(m);
			}else{
				$('#listjs').html(alertdiv('已经是第一页了','am-alert-danger')).show()
			}
	  });
	  $("#downpage").on('click',function(){
		  var n	= $('#novel-list');
		  var m = n.find("option:selected").next().val();
		  if(m){
			  n.val(m);
			  novel_list_ajax(m);
		  }else{
			  $('#listjs').html(alertdiv('已经是最后一页了','am-alert-danger')).show()
		  }
	  });
	  var options = {
			  target:'#gotext',
			  success:showResponse,
			  dataType:'json',
	  };
	  $('#doc-vld-msg').ajaxForm(options);
	  function novel_list_ajax(val){
		    $.ajax({
				 type:"get",
				 url:val,
				 dataType:"html",
				 timeout:"10000",
				 //data:{id:articleid,cid:chapterid},
				 beforeSend:function(){
					$('#listjs').html('<div class=\"am-progress am-progress-striped am-progress-sm am-active \">'+
							'<div class=\"am-progress-bar am-progress-bar-secondary\"  style=\"width: 100%\"></div>'+
							'</div>');
				 },
				 error:function(XMLHttpRequest){
					 $('#listjs').html(alertdiv('请求失败请重试')).show()
				 },
				 success:function(data){
				 	var shtml = "",
				 			Ljson = $.parseJSON(data);
					 	$.each(Ljson.List,function(i,k){
					    	shtml +='<li class="am-g"><a href="/read/'+k.Articleid+'/'+k.Chapterid+'" class="am-list-item-hd">'+k.Chaptername+'</a></li>';
						});
					 $('#listjs').html(alertdiv('当前页: '+$("#novel-list").find("option:selected").text())).show()
					 $('#novel-list-ajax').html(shtml);
					 $('#novel-list-page').html(Ljson.Page);
				 },
			  });
	  }
	  function nblackboard(val){
		  $('body,header').removeClass(function(){
		      return "am-header-default am-read-def am-read-eye am-read-eye2 am-read-off am-read-off2 am-read-zyg"
		    });
		  $("a[role='button'],button[id='bookmark']").removeClass(function(){
		      return "am-btn-secondary"
		    });

		  $("body").addClass(val);
		  if(val == 'am-read-def' ){
			  $("a[role='button'],button[id='bookmark']").addClass("am-btn-secondary");
			  $('header').addClass('am-header-default');
		  }
		  readcookie('mblackboard',val)
		  readse("#readbox",val)
	  }
	  function nsize(val){
		  $('#noveltext').removeClass(function(){
		      return "am-text-xs am-text-sm am-text-default am-text-lg am-text-xl am-text-xxl"
		    });
		   //删除size
		  $('#noveltext').addClass(val);
		  readcookie('msize',val)
		  readse("#readtxt",val)
	  }
	  function readcookie(name,val){
		  if($.cookie(name)){
			  $.cookie(name, val, {expires: 7, path: '/'});
		  }else{
			  $.cookie(name, val, {expires: 7, path: '/'});
		  }
	  }
	  function readse(name,val){
		  var $o =$(name).find('option[value="'+val+'"]');
		  $o.attr('selected', val);
	  }
	  function showResponse(responseText, statusText){
		  $('#alert-text-1').alert('close');
		  if(responseText.Status){
			  $('#alert-text').html(alertdiv(responseText.Info,'am-alert-success')).show();
			  jumurl(responseText.Url);
		  }else{
			  $('#alert-text').html(alertdiv(responseText.Info,'am-alert-danger')).show().delay(10000).fadeOut();

		  }
	  }
	  function alertdiv(text,color){
		  return '<div class=\"am-alert '+color+'\" id=\"alert-text-1\" data-am-alert>'+
		  '<button type=\"button\" class=\"am-close\">&times;</button><p>'+
		  text+
		  '</p></div>'
	  }
	  function jumurl(url){
		  window.location.href = url;
	  }
  });
})(jQuery);

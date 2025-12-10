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
  //简单表单验证
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
		  var $btn = $(this)
		  $btn.button('loading');
		  $.ajax({
			 type:"post",
			 url:"/wap/bookmark_add",
			 dataType:"json",
			 timeout:"3000",
			 data:{id:articleid,cid:chapterid},
			 success:function(data){
				 $btn.button('reset');
				 if(data.status!=0){
					 $btn.text(data.info);
				 }else{
					 $btn.text(reerror(data.info));
				 }
				 return false;
			 }
		  });
		  //setTimeout($('#bookmark').popover('close'),2000);
	  });
	  $('#novel-list').on('change', function() {
		    /*$('#listjs').html([
		      '选中项：<strong class="am-text-danger">',
		      [$(this).find('option').eq(this.selectedIndex).text()],
		      '</strong> 值：<strong class="am-text-warning">',
		      $(this).val(),
		      '</strong>'
		    ].join(''));*/
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
				 timeout:"3000",
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
					 
					 $('#listjs').html(alertdiv('当前页: '+$("#novel-list").find("option:selected").text())).show()
					 $('#novel-list-ajax').html(data);
					 
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
			  $.cookie(name, val);
		  }else{
			  $.cookie(name, val, {expires: 7, path: '/'});
		  }
	  }
	  function readse(name,val){
		  var $o =$(name).find('option[value="'+val+'"]');
		  $o.attr('selected', val);
	  }
	  function showResponse(responseText, statusText){
		  //alert('状态：'+statusText+'内容'+responseText.info);
/*		  $('#loginbtn').popover({
			    content: responseText.info,
		  }).popover('toggle')*/
		  $('#alert-text-1').alert('close');
		  if(responseText.status!=0){
			  $('#alert-text').html(alertdiv(responseText.info,'am-alert-success')).show();
			  jumurl(responseText.url);
		  }else{
			  $('#alert-text').html(alertdiv(reerror(responseText.info),'am-alert-danger')).show().delay(3000).fadeOut();

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
	  function reerror(re){
		  switch(re)
		  {
		  case -1:
			  return '用户名长度不合法';
			  break;
		  case -2:
			  return '用户名含有特殊字符';
			  break;
		  case -3:
			  return '密码长度不合法';
			  break;
		  case -4:
			  return '邮箱格式不正确';
			  break;
		  case -5:
			  return '用户名已被注册';
			  break;
		  case -10:
			  return '用户名或密码错误';
			  break;
		  case -11:
			  return  '没有此书';
			  break;
		  case -12:
			  return '没有此章';
			  break;
		  case -13:
			  return '用户名或密码错误';
			  break;
		  case -14:
			  return '用户名或密码错误';
			  break;
		  case -15:
			  return '用户名或密码错误';
			  break;
		  case -20:
			  return '已加签';
			  break;
		  case -21:
			  return '没有此书';
			  break;
		  case -22:
			  return '更签失败';
			  break;
		  case -23:
			  return '加签失败';
			  break;
		  default:
			  return '未知错误';
		  }
	  }

  });
})(jQuery);

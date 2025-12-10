$(function() {
	
	$('#user-book-list').find('.am-icon-close').add('#am-user-mark-del').
	on('click', function() {
		$('#my-confirm').modal({
			relatedTarget: this,
			onConfirm: function(options) {
				var $link = $(this.relatedTarget).prev("span").prev('a');
				//var delid = new Array();
				var delid = new Object();
				delid.bookid = new Array();
				if($link.length) {
					delid.bookid.push(Number($link.data('id')));
				} else {
					$("input[name='checkbox']:checked").each(function() {
						delid.bookid.push(Number($(this).val()));
					});
				}
				var $delBtn = $("#am-user-mark-del");
				$.ajax({
					type: "post",
					url: delUrl,
					dataType: "json",
					timeout: "10000",
					contentType : "application/json",
					data: JSON.stringify(delid),
					beforeSend: function() {
						$delBtn.button('loading');
					},
					error: function(xhr) {
						$delBtn.button('reset').text('请重试!');
					},
					success: function(data) {
						$delBtn.button('reset');
						alert(data.Info);
						if(data.Status){
							window.location.reload();
						}
						//window.location.reload();
						
					}
				});
			},
			onCancel: function() {
				return false;
			}
		});
	});
	$('#am-user-mark-move').click(function(){
		var $moveBtn	= $(this);
		var moveid = new Object();
		moveid.action =Number($('#mark-list-set').val())
		moveid.bookid = new Array();
		$("input[name='checkbox']:checked").each(function(){
            moveid.bookid.push(Number($(this).val()));
        });
		$.ajax({
			type: "post",
			url: moveUrl,
			contentType : "application/json",
			dataType: "json",
			timeout: "10000",
			data: JSON.stringify(moveid),
			beforeSend: function() {
				$moveBtn.button('loading');
			},
			error: function(xhr) {
				$moveBtn.button('reset').text('请重试');
			},
			success: function(data) {
				$moveBtn.button('reset');
				alert(data.Info);
				if(data.Status){
					window.location.reload();
				}

			}
		});
	});
	var verifyimg =$(".verify-imgl").attr("src");
	$(".vchange").click(function(){
		if(verifyimg.indexOf('?')>0){
			$(".verify-imgl").attr("src",verifyimg+'&reload='+Math.random());
		}else{
			$(".verify-imgl").attr("src",verifyimg+'?reload='+Math.random());
		}
	});
	$("#marklistclass").change(function(){
		//$(this).val();
		$(location).attr('href', bookCaseUrl+markclassorder+'/'+$(this).val());
	})
	$("#marklistorder").change(function(){
		//$(this).val();
		//alert(markclass);
		$(location).attr('href', bookCaseUrl+$(this).val()+'/'+markclass);
	})
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
    },
  });

  var options = {
		  beforeSubmit:showRequest,
		  success:showResponse,
		  dataType:'json',
		  timeout:5000,
  }
  $('#doc-vld-msg').ajaxForm(options);
  function showRequest(formData,jqForm,options){
	  if($('#doc-vld-msg').validator('validateForm').valid){
		  return true;
	  }else{
		  return false;
	  }
  }
  function showResponse(responseText,statusText){
	  var showid = $('#user_show');
	  //alert(statusText+responseText.info+responseText.url);
	  //$('#user_show').removeClass(function(){return 'am-alert-warning am-alert-success'});
	  
	  if(responseText.Status){
			showid.html(realert(responseText.Info,'am-alert-success'));
	  }else{
	  		showid.html(realert(responseText.Info,'am-alert-warning'));
	  }
	  setTimeout(function(){
	  	showid.html('');
	  },3000);
//	  $('#user_show').text(responseText.info);
//	  
//	  $('#user_show').show()
//	  
	  if(responseText.Url){
		  setTimeout(window.location.href=responseText.Url,3000);
	  }
  }
  function realert(reText,reState){
  	return '<div class="am-alert '+reState+'" data-am-alert><button type="button" class="am-close">&times;</button>'+reText+'</div>'
  }
});
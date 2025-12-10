(function($) {
    $('#getbaidu').on('click',function(){
    	searchStart(1);
    	//return false;
    });
    $('#spage').on('click', 'li', function() {
    	searchStart($(this).find("a").eq(0).attr("value"));
    	//return false;
    });
})(jQuery);
function getSearchData (data) {
	var shtml = ""
	$.each(data,function(i,k){
    	shtml +=assembleHtml(k);
	});
	$('#listcontent').html(shtml);
}
function searchStart(num){
	var num = arguments[0] ? arguments[0] : 1;
	var Search = new Object();
	Search.q = $("#baidusearch").val();
	Search.p = Number(num);
	$.ajax({
		url:"/soapi",
		type:"POST",
		timeout: 5000,
		dataType : "json",
		//jsonp : 'callback', 
		//jsonpCallback: 'handleResponse',
		data:JSON.stringify(Search),
		success:function(result){
			if(result.Status){
				if(result.nums == 0){
					//$('#searchHtml').html('<li>'+result.info+'</li>');
					alert(result.Info);
				}else{
					getSearchData(result.List);
					getSearchInfoData(result.Page,result.Pages);
				}
			}else{
				//$('#searchHtml').html('<li>'+result.info+'</li>');
				alert(result.Info);
			}
		
		},
		error:function(XMLHttpRequest, textStatus, errorThrown){
			//$('#searchHtml').html('<li>未知错误</li>');
			alert('未知错误');
			//console.log(XMLHttpRequest.status);
			//console.log(XMLHttpRequest.readyState);
			//console.log(textStatus);
		},
	});
}
function getSearchInfoData (num,nums) {
		var nownum = Number(num),
		max = 7,
		// totalnum = Math.ceil(Number(data.totalNum)/10),
		totalnum = Number(nums),
		cmax = Math.floor(max/2),
		tpage = nownum==1?'class="am-disabled"':'',
		dpage = nownum==totalnum?'class="am-disabled"':'';
		spage = '<li '+tpage+'><a id="page" value=1>&laquo;</a></li>'; 
	for($i=-cmax;$i<=cmax;$i++){
		if(nownum+$i>0&&nownum+$i<=totalnum){
			if(nownum+$i == nownum){
				spage+='<li class="am-active"><a id="opage" value='+Number(nownum+$i)+'>'+Number(nownum+$i)+'</a></li>';
			}else{
				spage+='<li><a id="opage" value='+Number(nownum+$i)+'>'+Number(nownum+$i)+'</a></li>';
			}
		}
	}
	spage +='<li '+dpage+'><a id="page" value='+totalnum+'>&raquo;</a></li>'; 
	$('#spage').html(spage);
   // return true;
}
function assembleHtml(data){
	var sort = "连载";
	if (data.Fullflag!=0){
		sort = "完结";
	}
    return '<li class="am-g am-list-item-desced am-list-item-thumbed am-list-item-thumb-bottom-left">'+
    '<div class="am-list-item-hd">'+
    '<a class="am-u-sm-8 am-text-truncate" href="/novel/'+data.Id+'" class="">'+data.Name+'</a>'+
    '<div class="am-u-sm-4 am-text-truncate am-text-right">'+data.Author+'</div>'+
    '</div>'+
    '<div class="am-u-sm-4 am-list-thumb">'+
    '<a href="/novel/'+data.Id+'" class="">'+
      '<img src="'+data.Imgurl+'" alt="'+data.Name+'"/>'+
    '</a>'+
    '</div>'+
	'<div class=" am-u-sm-8  am-list-main">'+
      '<div class="am-list-item-text">'+data.Intro+'</div>'+
	'</div>'+
	'</li>';
}
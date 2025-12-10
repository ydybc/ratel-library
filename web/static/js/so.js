(function($) {
	$('#startSearch').on('click', function() {
		listobtain($('#searchText').val(),1);
	});

	$('#searchPage').on("click","li",function(){
	    var p = $(this).find("a").eq(0).attr('value');
		listobtain($('#searchText').val(),Number(p));
	});
	var s = getQueryString("q");
	$('#searchText').val(s)
})(jQuery);

//组合结果
function assembleHtml(data){
	var sort = "连载";
	if (data.Fullflag!=0){
		sort = "完结";
	}
	return '<li>'+
	'<div class="am-g">'+
	'<a href="/book/'+data.Id+'"><img src="'+data.Imgurl+'" alt="'+data.Name+'" class="am-u-lg-3 am-radius" /></a>'+
	'<div class="am-u-lg-9">'+
	  		'<div class="am-text-truncate"><a href=/book/'+data.Id+'><h3>'+data.Name+'</h3></a></div>'+
		  	'<div class="am-text-truncate">作者：'+data.Author+'</div>'+
		  	'<div class="am-text-truncate">分类：'+data.Class+'</div>'+
		  	'<div class="am-text-truncate">状态：'+sort+'</div>'+
			'<div class="am-text-truncate">最新章节：<a href="/read/'+data.Id+'/'+data.Cid+'">'+data.Cname+'</a></div>'+
		  	'<div class="am-bookinfo-p line-clamp">'+data.Intro+'</div>'+
	'</div>'+
	'</div>'+
	'</li>';
}
function listobtain(query,page){
		var Search = new Object();
		Search.q = query;
		Search.p = page;
			$.ajax({
			type:"post",
			url:"/soapi",
			dataType:"json",
			timeout:"5000",
			data:JSON.stringify(Search),
			beforeSend:function(){
				
			},
			error:function(xhr){
				alert('请重试!');
			},
			success:function(data){
				if (data.Status){
					console.log(data)
					var shtml = "";
					$.each(data.List,function(i,k){
				    	shtml +=assembleHtml(k);
				    	//console.log(k);
					});
					$('#searchHtml').html(shtml);
					getSearchInfoData(data)
				}
			}
		});
}
function getSearchInfoData(data){
	var nownum = data.Page,
		max = 5,
		totalnum = data.Pages,
		cmax = Math.floor(max/2),
		tpage = nownum==1?'class="am-disabled"':'',
		dpage = nownum==totalnum?'class="am-disabled"':'';
		spage = '<li '+tpage+'><a id="page" value=1>&laquo;</a></li>'; 
	for($i=-cmax;$i<=cmax;$i++){
		if(nownum+$i>0&&nownum+$i<=totalnum){
			if(nownum+$i == nownum){
				spage+='<li class="am-active"><a id="page" value='+Number(nownum+$i)+'>'+Number(nownum+$i)+'</a></li>';
			}else{
				spage+='<li><a id="page" value='+Number(nownum+$i)+'>'+Number(nownum+$i)+'</a></li>';
			}
			
		}
	}
	spage +='<li '+dpage+'><a id="page" value='+totalnum+'>&raquo;</a></li>'; 
	$('#searchPage').html(spage);
}
function loading(){
	$('#searchHtml').html('<li><div class="am-progress am-progress-striped am-progress-sm am-active " style=" margin: 0;">'+
	  '<div class="am-progress-bar am-progress-bar-secondary"  style="width: 100%"></div></div></li>');
}
function getQueryString(name) {
	var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
	var r = window.location.search.substr(1).match(reg);
	if (r != null) return decodeURI(r[2]); return null;
}
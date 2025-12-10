(function($) {
	$('#startSearch').on('click',function(){
		loadScript();
		return false;
	});
    $('body').on('click', '#page', function() {
    	pageinit($(this).attr("value"));
    	return false;
    });
})(jQuery);
var cse,bds;
$(function(){ 
	bds = $("#searchText").val();
    if(bds){
    	loading();
		loadScript();
		return false;
    }
});

//初始化回调函数示例
function init () {
	bds = $("#searchText").val();
	if(!bds){
		alert('请输入您要搜索的信息');
		return false;
	}
	loading();
    cse = new BCse.Search("4836449536212096363");    //参数为您的API引擎ID，已自动填写，必需。
    cse.setResultType(2);
    cse.setPageNum(10);
    cse.getResult(bds, getSearchData,1);    //此方法获取搜索结果，参数1为搜索词，参数2为您获取到结果后想要执行的回调函数。
}
function pageinit (num) {
	loading();
    cse.getResult(bds, getSearchData,num);    //此方法获取搜索结果，参数1为搜索词，参数2为您获取到结果后想要执行的回调函数。
}
//处理结果回调函数示例
function getSearchData (data) {
	var shtml = ""
	$.each(data,function(i,k){
    	shtml +=assembleHtml(k);
    	//console.log(k);
	});
	$('#searchHtml').html(shtml);
	cse.getSearchInfo(bds,getSearchInfoData);
}
//组合结果
function assembleHtml(data){
	return '<li>'+
	'<div class="am-g">'+
	'<a href="'+data.linkUrl+'"><img src="'+data.image+'" alt="'+data.title+'" class="am-u-lg-3 am-radius" /></a>'+
	'<div class="am-u-lg-9">'+
	  		'<div class="am-text-truncate"><a href="'+data.linkUrl+'"><h3>'+data.title+'</h3></a></div>'+
		  	'<div class="am-text-truncate">作者：'+data.author+'</div>'+
		  	'<div class="am-text-truncate">分类：'+data.genre+'</div>'+
		  	'<div class="am-text-truncate">状态：'+data.updateStatus+'</div>'+
		  	'<div class="am-bookinfo-p line-clamp">'+data.abstract+'</div>'+
	'</div>'+
	'</div>'+
	'</li>';
}
function getSearchInfoData(data){
	var nownum = Number(data.curPage+1),
		max = 7,
		totalnum = Math.ceil(Number(data.totalNum)/10),
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
function loadScript () { 
    var script = document.createElement("script"); 
    script.type = "text/javascript";
    script.charset = "utf-8";
    script.src = "http://zhannei.baidu.com/api/customsearch/apiaccept?sid=4836449536212096363&v=2.0&callback=init";
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(script, s);
}
(function($) {
    $('#getbaidu').on('click',function(){
    	loadScript();
    	return false;
    });
    $('body').on('click', '#opage', function() {
    	initpage($(this).attr("value"));
    	return false;
    });
})(jQuery);
var cse,bds;
function display (data) {
    var listcon=""; 
    $.each(data,function(i,k){
    	listcon +=listcontent(k);
    });
    $('#listcontent').html(listcon);
    cse.getSearchInfo(bds,getinfo);
}
function getinfo (data) {
	var bdnum = Math.ceil(Number(data.totalNum)/10);
	var tpage,bpage,inunm = 0,bothnum = 2;
	if(data.curPage){
		inunm = Number(data.curPage);
	}
	var spage = '<li><a id="opage" value="1">&laquo;</a></li>';
	for(var i=bothnum; i>=1; i--){
		var tpage = inunm-i;
		if(tpage < 0){
			continue; 
		}else{
			spage += '<li><a id="opage" value="'+(tpage+1)+'">'+(tpage+1)+'</a></li>';
		}
	}
	spage += '<li class="am-active"><a id="opage" value="'+(inunm+1)+'">'+(inunm+1)+'</a></li>';
	for(var i=1;i<=bothnum;i++){
		var bpage = inunm+i;
		if(bpage >= bdnum){
			break;
		}else{
			spage += '<li><a id="opage" value="'+(bpage+1)+'">'+(bpage+1)+'</a></li>';
		}
	}
	spage+='<li><a id="opage" value="'+bdnum+'">&raquo;</a></li>';
	$('#spage').html(spage);
    return true;
}
function listcontent(data){
    return '<li class="am-g am-list-item-desced am-list-item-thumbed am-list-item-thumb-bottom-left">'+
    '<div class="am-list-item-hd">'+
    '<a class="am-u-sm-8 am-text-truncate" href="'+reurl(data.linkUrl)+'" class="">'+data.title+'</a>'+
    '<div class="am-u-sm-4 am-text-truncate am-text-right">'+data.author+'</div>'+
    '</div>'+
    '<div class="am-u-sm-4 am-list-thumb">'+
    '<a href="'+reurl(data.linkUrl)+'" class="">'+
      '<img src="'+data.image+'" alt="'+data.title+'"/>'+
    '</a>'+
    '</div>'+
	'<div class=" am-u-sm-8  am-list-main">'+
      '<div class="am-list-item-text">'+data.abstract+'</div>'+
	'</div>'+
	'</li>';
}
function reurl(url){
	return url.replace(/http:\/\/www.ziyouge.com\/zy\/(\d+)\/(\d+)\/index\.html/, "http://m.ziyouge.com/novel/$2");
}
function init () {
	bds = $("#baidusearch").val();
	if(!bds){
		alert('请输入您要搜索的信息');
		return false;
	}
	loading();
	var id = arguments[0] ? arguments[0] : 1; 
    cse = new BCse.Search("4836449536212096363");
    cse.setResultType(2);
    //cse.setPageNum(2);
    cse.getResult(bds, display, id);
}
function initpage(page){
	if(!bds){
		alert('请输入您要搜索的信息');
		return false;
	}
	loading();
	cse.setResultType(2);
	cse.getResult(bds, display, page);
}
function loading(){
	$('#listcontent').html(	'<div class="am-progress am-progress-striped am-progress-sm am-active ">'+
			 '<div class="am-progress-bar am-progress-bar-secondary"  style="width: 100%"></div></div>');

}
function loadScript () { 
    var script = document.createElement("script"); 
    script.type = "text/javascript";
    script.charset = "utf-8";
    script.src = "http://zhannei.baidu.com/api/customsearch/apiaccept?sid=4836449536212096363&v=2.0&callback=init";
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(script, s);
}
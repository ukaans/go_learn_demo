$(function(){
  baseApp.init();
})

var baseApp = {
  init: function(){
      this.initAside();
      this.confirmDelete();   
      this.resizeIframe();
      this.changeStatus();
  },
  initAside: function(){
      $('.aside h4').click(function(){
          $(this).siblings('ul').slideToggle();
      });
  },
  resizeIframe: function(){					
      $("#rightMain").height($(window).height() - 80);
  },
  // 删除提示
  confirmDelete: function(){
      $(".delete").click(function(e){
          if (!confirm("您确定要删除吗？")) {
              e.preventDefault();
          }
      });
  },
  // 点击图标修改状态
  changeStatus:function(){
    $(".chStatus").click(function(){
        var id=$(this).attr("data-id")
        var table=$(this).attr("data-table")
        var field=$(this).attr("data-field")
        var el =$(this)
        $.get("/admin/changeStatus",{id:id,table:table,field:field},function(response){				
            if(response.success){
                if (el.attr("src").indexOf("yes")!=-1){
                    el.attr("src","/static/admin/images/no.gif")
                }else{
                    el.attr("src","/static/admin/images/yes.gif")
                }
            }
        })
    })
}
}
(function($){

    var app={
        init:function(){
    
            this.initSwiper();

            this.initNavSlide();

            this.initProductContentTab();

            this.initProductContentColor();

            this.initProductContentAttr();
        },
        initSwiper:function(){    
            new Swiper('.swiper-container', {
                loop : true,
                navigation: {
                  nextEl: '.swiper-button-next',
                  prevEl: '.swiper-button-prev'                 
                },
                pagination: {
                    el: '.swiper-pagination',
                    clickable :true
                }
                
            });
        },
        initNavSlide:function(){
             $("#nav_list>li").hover(function(){

                $(this).find('.children-list').show();
             },function(){
                $(this).find('.children-list').hide(); 
             })          

        },
        initProductContentTab: function(){
            // 默认激活第一个
            $('.detail_info .detail_info_item:first').addClass('active');
            $('.detail_list li:first').addClass('active');
        
            // Tab 点击切换
            $('.detail_list li').click(function () {
                var index = $(this).index();
                $(this).addClass('active').siblings().removeClass('active');
                $('.detail_info .detail_info_item').removeClass('active').eq(index).addClass('active');
            });
        },
        initProductContentColor:function(){
            var _that=this;
            $("#color_list .banben").first().addClass("active");
            $("#color_name").html($("#color_list .active .yanse").html())
            $("#color_list .banben").click(function(){
                $(this).addClass("active").siblings().removeClass("active");                
                $("#color_name").html($("#color_list .active .yanse").html())
                var goods_id=$(this).attr("goods_id")
                var color_id=$(this).attr("color_id")

                $.get("/product/getImgList",{"goods_id":goods_id,"color_id":color_id},function(response){
                    console.log(response)
                    if(response.success==true){
                        var swiperStr=""
                        for (var i = 0; i < response.result.length; i++) {
                            swiperStr += '<div class="swiper-slide"><img src="' + response.result[i].img_url + '"> </div>';                            
                        }
                        $("#item_focus").html(swiperStr)
                        _that.initSwiper()
                    }
                })
            })
        },
        initProductContentAttr: function() {
            var _that = this;
            
            $(".xzbb:not(#color_list) .banben").each(function() {
                var $banben = $(this);
                
                if ($banben.index() === 0) {
                    $banben.addClass("active");
                }
                
                $banben.click(function() {
                    $(this).addClass("active").siblings(".banben").removeClass("active");
                    _that.updateSelectedInfo();
                });
            });
        
            // 修复：定义缺失的方法
            this.updateSelectedInfo = function() {
                // TODO: 这里可以实现右侧价格、版本、颜色同步更新
                console.log("[updateSelectedInfo] 已执行");
                
                // 示例：同步颜色名称（颜色部分已经在 initProductContentColor 里处理了）
                var colorName = $("#color_list .active .yanse").html() || "";
                $("#color_name").html(colorName);
            };
        
            // 初始化调用
            this.updateSelectedInfo();
        },
           
    }   
    
    $(function(){
    
    
        app.init();
    })

    

})($)

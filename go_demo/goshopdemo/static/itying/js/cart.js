(function($){
    var app={
        init:function(){             
            this.changeCartNum();       
            this.deleteConfirm();
            this.initCheckBox();
            this.isCheckedAll();
            this.initChekOut();
        },

        // 新增：统一更新购物车数量显示
        updateCartCount: function(total, selected) {
            if(total !== undefined && selected !== undefined) {
                // 更新 "共X件商品，已选择Y件"
                $(".tishi li").eq(2).find("span").eq(0).text(total);   // 共X件
                $(".tishi li").eq(2).find("span").eq(1).text(selected); // 已选择Y件
            }
        },

        deleteConfirm:function(){
            $('.delete').click(function(){    
                var flag=confirm('您确定要删除吗?');    
                return flag;    
            })    
        },     

        initCheckBox: function(){
            var _that = this;

            //全选按钮点击
            $("#checkAll").click(function() {               
                if (this.checked) {
                    $(":checkbox").prop("checked", true);
                    $.get('/cart/changeAllCart?flag=1',function(response){                 
                        if(response.success){
                            $("#allPrice").html(response.allPrice + "元");
                            _that.updateCartCount(response.totalCount, response.selectedCount);
                        }
                    })
                } else {
                    $(":checkbox").prop("checked", false);      
                    $.get('/cart/changeAllCart?flag=0',function(response){                 
                        if(response.success){
                            $("#allPrice").html(response.allPrice + "元");
                            _that.updateCartCount(response.totalCount, response.selectedCount);
                        }
                    })                           
                }               
            });    

            //点击单个选择框
            $(".cart_list :checkbox").click(function() {                         
                _that.isCheckedAll();

                var goods_id = $(this).attr("goods_id");
                var goods_color = $(this).attr("goods_color");

                $.get('/cart/changeOneCart?goods_id='+goods_id+'&goods_color='+goods_color, function(response){                 
                    if(response.success){
                        $("#allPrice").html(response.allPrice + "元");
                        _that.updateCartCount(response.totalCount, response.selectedCount);
                    }
                });
            });   
        },

        //判断全选是否选中
        isCheckedAll: function(){             
            var allNum = $(".cart_list :checkbox").size();
            var checkedNum = 0;  

            $(".cart_list :checkbox").each(function () {  
                if($(this).prop("checked") == true){
                    checkedNum++;
                }
            });

            if(allNum == checkedNum){
                $("#checkAll").prop("checked", true);
            } else {
                $("#checkAll").prop("checked", false);
            }
        }, 

        changeCartNum: function(){
            var _that = this;

            // 减少数量
            $('.decCart').click(function(){
                var goods_id = $(this).attr("goods_id");
                var goods_color = $(this).attr("goods_color");
                var currentBtn = this;

                $.get('/cart/decCart?goods_id='+goods_id+'&goods_color='+goods_color, function(response){
                    if(response.success){
                        $("#allPrice").html(response.allPrice + "元");
                        $(currentBtn).siblings(".input_center").find("input").val(response.num);
                        $(currentBtn).parent().parent().siblings(".totalPrice").html(response.currentPrice + "元");
                        
                        // 更新顶部数量
                        _that.updateCartCount(response.totalCount, response.selectedCount);
                    }
                });
            });

            // 增加数量
            $('.incCart').click(function(){
                var goods_id = $(this).attr("goods_id");
                var goods_color = $(this).attr("goods_color");
                var currentBtn = this;

                $.get('/cart/incCart?goods_id='+goods_id+'&goods_color='+goods_color, function(response){
                    if(response.success){
                        $("#allPrice").html(response.allPrice + "元");
                        $(currentBtn).siblings(".input_center").find("input").val(response.num);
                        $(currentBtn).parent().parent().siblings(".totalPrice").html(response.currentPrice + "元");
                        
                        // 更新顶部数量
                        _that.updateCartCount(response.totalCount, response.selectedCount);
                    }
                });
            });
        },

        initChekOut: function(){
            $("#checkout").click(function(){	
                var allPrice = parseFloat($("#allPrice").html());	
                if(allPrice == 0){
                    alert('购物车没有选中去结算的商品');
                } else {
                    location.href = "/buy/checkout";
                }	
            });
        },
    };

    $(function(){
        app.init();
    });    
})($);
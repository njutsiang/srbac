$(function(){
    $('.main-auto-search-form').each(function(){
        var that = $(this);
        that.find('select').change(function(){
            that.submit();
        });
    });
    $('.delete-form').each(function(){
        var that = $(this);
        that.find('.btn').click(function(){
            if (confirm("确定要删除该数据？")) {
                that.submit();
            }
        });
    });
    $('.form-post').each(function(){
        var that = $(this);
        that.find('.form-post-btn').click(function(){
            if (confirm($(this).attr('data-message'))) {
                that.submit();
            }
        });
    });
});
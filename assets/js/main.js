$(function(){
    $('.main-auto-search-form').each(function(){
        var that = $(this);
        that.find('select').change(function(){
            that.submit();
        });
    });
});
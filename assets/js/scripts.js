(function ($) {
    "use strict";

    $("input[type='number']").attr('min', 0);

    $('#cost, #quantity').change(function () {
        var price = parseFloat($('#cost').val());
        var qty = parseFloat($('#quantity').val());
        $('#subTotal').val(price * qty);
    });

    $('#unit_rate_grn, #qty_received').change(function () {
        var price = parseFloat($('#unit_rate_grn').val());
        var qty = parseFloat($('#qty_received').val());
        $('#amount_grn').val(price * qty);
    });

    $('#unit_rate_gin, #qty_issued').change(function () {
        var price = parseFloat($('#unit_rate_gin').val());
        var qty = parseFloat($('#qty_issued').val());
        $('#amount_gin').val(price * qty);
    });

    $('#issue_unit_rate, #issue_qty').change(function () {
        var price = parseFloat($('#issue_unit_rate').val());
        var qty = parseFloat($('#issue_qty').val());
        $('#issue_amount').val(price * qty);
    });

    $('#balance_unit_rate, #balance_qty').change(function () {
        var price = parseFloat($('#balance_unit_rate').val());
        var qty = parseFloat($('#balance_qty').val());
        $('#balance_amount').val(price * qty);
    });

    $('#rcpt_unit_rate, #rcpt_qty').change(function () {
        var price = parseFloat($('#rcpt_unit_rate').val());
        var qty = parseFloat($('#rcpt_qty').val());
        $('#rcpt_amount').val(price * qty);
    });

    $('#issue_qty, #rcpt_qty').change(function () {
        var price = parseFloat($('#issue_qty').val());
        var qty = parseFloat($('#rcpt_qty').val());
        $('#balance_qty').val(qty - price);
    });

    $(document).ready(function () {
        $("#searchInput").on("keyup", function () {
            var value = $(this).val().toLowerCase();
            $("#searchTable tr").filter(function () {
                $(this).toggle($(this).text().toLowerCase().indexOf(value) > -1)
            });
        });
    });

    $(document).ready(function () {
        let max = 99999;
        let min = 1000;
        $("#barcodeItem").val("DIT" + "-" + "0" + (Math.floor(Math.random() * (max - min + 1) + min)));
    });

    $(document).ready(function () {
        let max = 9999;
        let min = 1;
        $("#indexSerialNo").val("0" + (Math.floor(Math.random() * (max - min + 1) + min)));
        $('#indexSerialNo').attr('min', 0);
    });

    $('.form-control').each(function () {

        var default_value = this.value;

        $(this).focus(function () {
            if (this.value == default_value) {
                this.value = '';
            }
        });

        $(this).blur(function () {
            if (this.value == '') {
                this.value = default_value;
            }
        });

    });

    $(document).ready(function () {
        var url = window.location;

        // Will only work if string in href matches with location
        $('ul.nav-content a[href="' + url + '"]').parent().addClass('active');

        // Will also work for relative and absolute hrefs
        $('ul.nav-content a').filter(function () {
            return this.href == url;
        }).parent().addClass('active').parent().parent().addClass('active');
    });

    $(function () {
        var dtToday = new Date();

        var month = dtToday.getMonth() + 1;
        var day = dtToday.getDate();
        var year = dtToday.getFullYear();
        if (month < 10)
            month = '0' + month.toString();
        if (day < 10)
            day = '0' + day.toString();

        var minDate = year + '-' + month + '-' + day;
        $('.date-control').attr('min', minDate).val(minDate);
    });

    $("input[name='date_taken']").change(function () {
        $("input[name='expect_return']").attr('min', $(this).val()).val($(this).val());
    });

    /* 
    ------------------------------------------------
    Sidebar open close animated humberger icon
    ------------------------------------------------*/

    $(".hamburger").on('click', function () {
        $(this).toggleClass("is-active");
    });


    /*  
    -------------------
    List item active
    -------------------*/
    $('.header li, .sidebar li').on('click', function () {
        $(".header li.active, .sidebar li.active").removeClass("active");
        $(this).addClass('active');
    });

    $(".header li").on("click", function (event) {
        event.stopPropagation();
    });

    $(document).on("click", function () {
        $(".header li").removeClass("active");

    });

    /*  
    -----------------
    Chat Sidebar
    ---------------------*/
    var open = false;

    var openSidebar = function () {
        $('.chat-sidebar').addClass('is-active');
        $('.chat-sidebar-icon').addClass('is-active');
        open = true;
    }
    var closeSidebar = function () {
        $('.chat-sidebar').removeClass('is-active');
        $('.chat-sidebar-icon').removeClass('is-active');
        open = false;
    }

    $('.chat-sidebar-icon').on('click', function (event) {
        event.stopPropagation();
        var toggle = open ? closeSidebar : openSidebar;
        toggle();
    });


    /*  Auto date in footer
    --------------------------------------*/

    document.getElementById("date-time").innerHTML = Date();


    /* TO DO LIST 
    --------------------*/
    $(".tdl-new").on('keypress', function (e) {
        var code = (e.keyCode ? e.keyCode : e.which);
        if (code == 13) {
            var v = $(this).val();
            var s = v.replace(/ +?/g, '');
            if (s == "") {
                return false;
            } else {
                $(".tdl-content ul").append("<li><label><input type='checkbox'><i></i><span>" + v + "</span><a href='#' class='ti-close'></a></label></li>");
                $(this).val("");
            }
        }
    });

    $(".tdl-content a").on("click", function () {
        var _li = $(this).parent().parent("li");
        _li.addClass("remove").stop().delay(100).slideUp("fast", function () {
            _li.remove();
        });
        return false;
    });

    // for dynamically created a tags
    $(".tdl-content").on('click', "a", function () {
        var _li = $(this).parent().parent("li");
        _li.addClass("remove").stop().delay(100).slideUp("fast", function () {
            _li.remove();
        });
        return false;
    });

    /*  Chat Sidebar User custom Search
    ---------------------------------------*/

    $('[data-search]').on('keyup', function () {
        var searchVal = $(this).val();
        var filterItems = $('[data-filter-item]');

        if (searchVal != '') {
            filterItems.addClass('hidden');
            $('[data-filter-item][data-filter-name*="' + searchVal.toLowerCase() + '"]').removeClass('hidden');
        } else {
            filterItems.removeClass('hidden');
        }
    });

    /*  Checkbox all
    ---------------------------------------*/

    $("#checkAll").change(function () {
        $("input:checkbox").prop('checked', $(this).prop("checked"));
    });


    /*  Vertical Carousel
    ---------------------------*/

    $('#verticalCarousel').carousel({
        interval: 2000
    })

    $(window).bind("resize", function () {
        console.log($(this).width())
        if ($(this).width() < 680) {
            $('.logo').addClass('hidden')
            $('.sidebar').removeClass('sidebar-shrink')
            $('.sidebar').removeClass('sidebar-shrink, sidebar-gestures')
        }
    }).trigger('resize');

    /*  Search
    ------------*/
    $('a[href="#search"]').on('click', function (event) {
        event.preventDefault();
        $('#search').addClass('open');
        $('#search > form > input[type="search"]').focus();
    });

    $('#search, #search button.close').on('click keyup', function (event) {
        if (event.target == this || event.target.className == 'close' || event.keyCode == 27) {
            $(this).removeClass('open');
        }
    });


    /*  pace Loader
    -------------*/

    paceOptions = {
        elements: true
    };

})(jQuery);
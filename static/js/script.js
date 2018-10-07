var digitMap = { "0": "۰", "1": "۱", "2": "۲", "3": "۳", "4": "۴", "5": "۵", "6": "۶", "7": "۷", "8": "۸", "9": "۹" };

function toFarsiNumber(n) {
    n = n.replace(/\d/g, function (d) {
        return digitMap["" + d];
    });
    return n;
}

var load_unix_time = Math.round(new Date().getTime() / 1000);

$(document).ready(function () {
    $(window).on('scroll', function (event) {
        $('div[data-spy]').each(function () {
            var top = parseInt($(this).data('offset-top'));
            var scrollValue = $(window).scrollTop();
            if (scrollValue > top)
                $(this).addClass('affixed');
            else
                $(this).removeClass('affixed');
        });
    });

    if ($("#particles-js").length > 0) {
        particlesJS.load("particles-js", '/static/data/particles-background.json?ver=' + load_unix_time);
    }

    $('.team-members').owlCarousel(
        {
            loop: true,
            margin: 10,
            responsiveClass: true,
            nav: true,
            autoplay: true,
            autoplayTimeout: 3000,
            smartSpeed: 700,
            items: 1,
        }
    );

    setTimeout(function () {
        $('.preloader').fadeOut(500);
    }, 500);

    $(".btn-target").each(function () {
        $(this).click(function () {
            var target = $(this).data('target');
            $.when($(".preloader").fadeIn(500)).done(function () {
                window.location = target;
            });
        });
    });
});

$(function () {
    $('[data-toggle="tooltip"]').tooltip()
});
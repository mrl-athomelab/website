$(document).ready(function () {
    $(".alert").each(function () {
        if ($(this).hasClass("hidden"))
            $(this).hide();
    });

    var administrators = $(".page-admin-administrators");
    administrators.find(".admin-table .remove-button").click(function () {
        var btn = $(this);
        btn.attr('disabled', 'disabled');
        var id = btn.data('id');
        $.ajax({
            type: "DELETE",
            url: "/admin/administrators/rest/" + id,
        }).done(function (data) {
            toastr.success(data.message, data.result.message, { timeOut: 3000 });
        }).fail(function (err) {
            console.log(err);
            toastr.error('error on submiting, please check console log !', 'error !', { timeOut: 3000 });
        }).always(function () {
            btn.removeAttr('disabled');
        });
    });
    administrators.find(".new-admin button").click(function () {
        var data = {
            "username": administrators.find("input[name='username']").val(),
            "password": administrators.find("input[name='password']").val(),
            "csrf_token": administrators.find("input[name='csrf_token']").val(),
        }
        var btn = $(this);
        btn.attr('disabled', 'disabled');

        $.ajax({
            type: "POST",
            url: "/admin/administrators/rest",
            data: JSON.stringify(data)
        }).done(function (data) {
            toastr.success(data.message, data.result.message, { timeOut: 3000 });
            location.reload();
        }).fail(function (err) {
            console.log(err);
            toastr.error('error on submiting, please check console log !', 'error !', { timeOut: 3000 });
        }).always(function () {
            btn.removeAttr('disabled');
        });
    });
});
$(document).ready(function () {
    $(".alert").each(function () {
        if ($(this).hasClass("hidden"))
            $(this).hide();
    });

    administrators_page();
    members_page();
});

function members_page() {
    var members = $(".page-admin-members");
    members.find(".new-member button").click(function () {
        var data = {
            "firstname": members.find("input[name='firstname']").val(),
            "lastname": members.find("input[name='lastname']").val(),
            "biography": members.find("textarea[name='biography']").val(),
            "socialmedialink": members.find("input[name='socialmedialink']").val(),
            "socialmediatype": members.find("select[name='socialmediatype']").val(),
        }
        console.log(JSON.stringify(data));
        // var btn = $(this);
        // btn.attr('disabled', 'disabled');

        // $.ajax({
        //     type: "POST",
        //     url: "/admin/members/rest",
        //     data: JSON.stringify(data)
        // }).done(function (data) {
        //     toastr.success(data.message, data.result.message, { timeOut: 3000 });
        //     location.reload();
        // }).fail(function (err) {
        //     console.log(err);
        //     toastr.error('error on submiting, please check console log !', 'error !', { timeOut: 3000 });
        // }).always(function () {
        //     btn.removeAttr('disabled');
        // });
    });
}

function administrators_page() {
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
}
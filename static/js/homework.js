$(document).ready(function () {
    $('.btn-link').unbind('click');
    $('body').on('click', 'button.btn-link', function(){
        $("#descriptionTitle").text($(this).text());
        $(".modal-body").html(descArr[$(this).attr('data-iter')]);
        $('#description').modal("show");
    });
});
var pages = [];
$("ul").each(function() {
    pages.push($(this).attr("id"));
});
var descArr = [];
var iter = 0;
for (var k = 0; k < pages.length; k++) {
    (function(k) {
        $.ajax({
            type: "GET",
            url: "/homework/assignments",
            data: {
                classnum: pages[k]
            },
            dataType: "json"
        }).done(function(data) {
            displayAssignment(data, pages[k]);
        });
    })(k);
}
function displayAssignment (data, id) {
    for (var i = 0; i < data.length; i++) {
        if (data[i].Category.trim() != "") {

            var period = pages.indexOf(id) + 1;
            if (period > 2) {
                period++
            }

            if (data[i].Category.indexOf(period) == -1) {
                continue;
            }
        }
        descArr.push(data[i].Desc);
        $("#" + id).append("<li><button class='btn-link' data-iter='" + iter + "'>" + data[i].Name + "</button><small>" + data[i].Date + "</small></li>");
        iter++
    }
}

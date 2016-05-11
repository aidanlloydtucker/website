$("#save").click(function(){
    var classlist = [];
    var errNotify = false;
    $("select").each(function() {
        if (!$(this).find("option:selected").attr("disabled")) {
            classlist.push($(this).find("option:selected").val());
        } else {
            errNotify = true;
        }
    });
    if (errNotify) {
        notification(4, "Error", "You Must Select A Class");
        return;
    }
    $.ajax({
        type: "PUT",
        url: "classes",
        data: {
            classlist: JSON.stringify(classlist)
        }
    }).done(function(data) {
        window.location.href = "/homework";
    });
});

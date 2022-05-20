$(document).ready(function() {
    $(".tablesorter").tablesorter();

    scrollAttendanceLists();

    $(".clickable").on('click', function() {
        if ($(this).next().is("div"))
            $(this).next().toggle()
        $(this).toggleClass("collapse")
    })

    $(".button.clear").on('click', function() {
        $(this).prev().children().prop("selected", false);
        return false;
    });
});


function scrollAttendanceLists() {
    $(".attendance .dates").scrollLeft(10000);
}


function clearSelect() {
    console.log($(this).prev());
}

function keep(element, which, a_id, b_id) {
    $.get("/import/duplicates/keep", {
        keep: which,
        a: a_id,
        b: b_id
    });
    $(element).parents(".compare_duplicates").remove();
}

function add_to_program(element, program, classType, id, fname, row, enroll) {
    $.get("/import/attendance", {
        stage: 'add_student',
        id: id,
        program: program,
        class_type: classType,
        filename: fname,
        row: row,
        enroll: enroll
    });
    $(element).parents(".compare_duplicates").remove();
}


function search(input, table) {
    var q = $(input).val().split(" ")
    var query = ""
    for (var i = 0; i < q.length; i++) {
        if (q[i].length > 0)
            query += ":contains(" + q[i] + ")"
    }
    $("." + table + " .row, ." + table + " tbody tr").hide()
    $("." + table + " .row" + query + ", ." + table + " tbody tr" + query).show()
}

function attended(box, id) {
    var checked = $(box).hasClass("checked")
    $.post("", {
        stage: 'student_status',
        id: id,
        attended: checked
    });
    // Uncheck the excused box
    if (checked) {
        var att = $(box).parent("td").next().children(0)
        att.removeClass("checked")
    }
}

function excused(box, id) {
    var checked = $(box).hasClass("checked")
    $.post("", {
        stage: 'student_status',
        id: id,
        excused: checked
    });
    //Uncheck the attended box
    if (checked) {
        var att = $(box).parent("td").prev().children(0)
        att.removeClass("checked")
    }
}

function maths(box, id) {
    var checked = $(box).hasClass("checked")
    $.post("", {
        stage: 'mathematics',
        id: id,
        math: checked
    });
    // Uncheck the excused box
    if (checked) {
        var att = $(box).parent("td").next().children(0)
        att.removeClass("checked")
    }
}

function physics(box, id) {
    var checked = $(box).hasClass("checked")
    $.post("", {
        stage: 'physical_sciences',
        id: id,
        phys: checked
    });
    // Uncheck the excused box
    if (checked) {
        var att = $(box).parent("td").next().children(0)
        att.removeClass("checked")
    }
}

function homework(box, id) {
    var checked = $(box).hasClass("checked")
        //Check if its yes or no
    var doneBox = $(box).hasClass("green")
    $.post("", {
        stage: 'homework',
        id: id,
        doneBox: doneBox,
        checked: checked
    });

    //Uncheck the other box
    if (checked) {
        var att = $(box).prev("div.custombox")
        att.removeClass("checked")
        att = $(box).next("div.custombox")
        att.removeClass("checked")
    }
}


function removeStudent(element, studentID) {
    var tr = $(element).parents("tr").remove();
    $.post("", {
        stage: 'remove',
        id: studentID
    })
}

// Make searching case insensitive
$.expr[":"].contains = $.expr.createPseudo(function(arg) {
    return function(elem) {
        return $(elem).text().toUpperCase().indexOf(arg.toUpperCase()) >= 0;
    };
});
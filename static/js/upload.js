
var RECEIVE_FILE_DROP = true;
function step1_drop(e) {
  e.preventDefault();
  e.stopPropagation();
  if (!RECEIVE_FILE_DROP) return;

  console.log(e);
  if (e.originalEvent.dataTransfer && e.originalEvent.dataTransfer.files.length==1){
    $("#selectedfile").prop("files", e.originalEvent.dataTransfer.files);
  }
}

function step1_filechanged() {
  if ($("#selectedfile").prop("files").length != 1) {
    $("#step1-error").show().text("Invalid selection. Please try again.");
    $("#step1 .upload-info").hide();
    $("#step1-fileselected").hide();
    return;
  }

  var file = $("#selectedfile").prop("files")[0];
  $(".upload-info .filename").text(file.name);
  $(".upload-info .filesize").text(file.size + " bytes");
  $(".upload-info .filetype").text(file.type);

  $("#step1-error").hide();
  $("#step1 .upload-info").show();
  $("#step1-fileselected").show();
}

function step1_go() {
  if ($("#selectedfile").prop("files").length != 1) {
    $("#step1-error").show().text("Invalid selection. Please try again.");
    $("#step1 .upload-info").hide();
    $("#step1-fileselected").hide();
    return;
  }

  var file = $("#selectedfile").prop("files")[0];
  $("#step1").fadeOut();

  $.ajax({
    method: "POST",
    url: "/api/new_upload",
    headers: {
      "X-FileBounce-Filename": file.name,
      "X-FileBounce-Content-Type": file.type,
      "X-FileBounce-Content-Length": file.size,
      "X-FileBounce-Token": "not implemented",
    },
    dataType: "text",
    success: function(data){
      $("#downloadid").val(data);
      var dlHref = window.location.origin + "/d/" + data;
      $("a#dllink").attr("href", dlHref).text(dlHref);
      $("#step2").fadeIn();
    }
  });

}

function step2_confirm() {
  $("#dl-instructions").fadeOut();
  $("#upload-trigger").fadeIn();
}

var CHUNK_SIZE = 20000; // 20 KB, arbitrary

function send_chunk(file, start, url) {
  var data = file.slice(start, start+CHUNK_SIZE);
  var reader = new FileReader();

  reader.onload = function() {
    var dataUrl = reader.result;
    var base64 = dataUrl.split(',')[1];
    $.ajax({
      method: "POST",
      url: url,
      data: base64,
      processData: false,
      success: function() {
        var new_start = start + CHUNK_SIZE;
        if (new_start > file.size) {
          return;
        }
        send_chunk(file, new_start, url);
      }
    });
  };
  reader.readAsDataURL(data);
}


function do_upload() {
  var dlId = $("#downloadid").val();
  var uploadUrl = window.location.origin + "/api/upload/" + dlId;
  var file = $("#selectedfile").prop("files")[0];
  send_chunk(file, 0, uploadUrl);
}

$(document).ready(function() {
  $(document).on('dragover', function(e) {
    e.preventDefault();
    e.stopPropagation();
  });
  $(document).on('dragenter', function(e) {
    e.preventDefault();
    e.stopPropagation();
  });
  $(document).on("drop", step1_drop)


  $("#step1 .file-pick").click(function(){
    $("#selectedfile").trigger("click");
  });

  $("#selectedfile").on("change", step1_filechanged);
  $("#step1-go").click(step1_go);
  $("#step2-confirm").click(step2_confirm);
  $("#upload-trigger").click(do_upload);
});

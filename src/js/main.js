$(function() {
    /* Most of this courtesy of http://paulrouget.com/miniuploader/ */
    window.ondragover = function(e) {
        e.preventDefault();
    }
    window.ondrop = function(e) {
        e.preventDefault();
        upload(e.dataTransfer.files[0]);
    }
});

function upload(file) {
    /* Is the file an image? */
    if (!file || !file.type.match(/image.*/)) {
        alert("That's not an image file.");
        return;
    }

    /* It is! */
    $(".upload-box").toggleClass("uploading");

    /* Lets build a FormData object*/
    var fd = new FormData(); // I wrote about it: https://hacks.mozilla.org/2011/01/how-to-develop-a-html5-image-uploader/
    fd.append("image", file); // Append the file
    var xhr = new XMLHttpRequest(); // Create the XHR (Cross-Domain XHR FTW!!!) Thank you sooooo much imgur.com
    xhr.open("POST", "https://api.imgur.com/3/image.json"); // Boooom!
    xhr.onload = function() {
        // Big win!
        var fullSize = JSON.parse(xhr.responseText).data.link;
        var imageUrl = fullSize.replace(".jpg", "h.jpg");
        apiAddSchmoopy(imageUrl);
    }
    xhr.onerror = function() {
        $(".upload-box").toggleClass("uploading");
        alert("Ah, shit. Upload failed. Sorry. Try again.");
    }

    xhr.setRequestHeader('Authorization', 'Client-ID 28aaa2e823b03b1'); // Get your own key http://api.imgur.com/

    /* And now, we send the formdata */
    xhr.send(fd);
}

function apiRemoveSchmoopy(imageUrl) {
    // TODO
}

function apiAddSchmoopy(imageUrl) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/schmoopy/add");
    /* use the large thumbnail */
    xhr.onload = function() {
        $(".upload-box").toggleClass("uploading");
        addSchmoopy(imageUrl);
    }
    xhr.onerror = function() {
        $(".upload-box").toggleClass("uploading");
        alert("Ah, shit. Couldn't save it. Sorry. Try again.");
    }

    /* And now, we send the formdata */
    var fd = new FormData();
    fd.append("name", g_name);
    fd.append("imageUrl", imageUrl);
    xhr.send(fd);
}

function addSchmoopy(imageUrl) {
    var div = document.createElement("div");
    $(div).addClass("schmoopy");

    var a = document.createElement("a");
    a.href = imageUrl;
    a.target = "_new";

    var img = document.createElement("img");
    img.src = imageUrl;

    $(img).appendTo(a);
    $(a).appendTo(div);
    $(".schmoopys").append(div);


    function queueAnimation() {
        var windowWidth = $(window).width();
        var windowHeight = $(window).height();

        var imgWidth = $(img).width();
        var imgHeight = $(img).height();

        var newLeft = Math.random()*(windowWidth-imgWidth);
        var newTop = Math.random()*(windowHeight-imgHeight);

        var newTime = Math.random()*(15000-3000)+3000;

        $(div).animate({
            left: newLeft,
            top: newTop,
        }, newTime, queueAnimation);
    }
    setTimeout(queueAnimation, 2000);
}

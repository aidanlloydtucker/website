var socket = io();
var username;
Notification.requestPermission();

var messages = [];
var numPresses = 0;

function setName() {
    user = $('#name-input').val().trim();
    $('#name-input').val('');
    if (!user) {
        return false;
    }
    socket.emit('chat nick', user);
}
function sendMsg() {
    var msg = $('#msg-input').val().trim();
    $('#msg-input').val('');
    if (!msg) {
        return false;
    }
    messages.push(msg);
    numPresses = messages.length;
    if (msg.charAt(0) === "/") {
        socket.emit('chat command', msg);
        if (msg.substr(0, 3) === "/dm") {
            var touser = msg.split(" ")[1];
            var dmArr = msg.split(" ");
            dmArr.splice(0, 2);
            var chatMsg = dmArr.join(" ");
            $('#messages').append('<p><span class="text-muted">DM to ' + touser + ':</span> ' + chatMsg + '</p>');
        }
        return;
    }
    socket.emit('chat message', msg);
    $('#messages').append('<p><span class="text-muted">' + username + ':</span> ' + msg + '</p>');
}

$("#img-input").change(function() {
    var file = this.files[0];
    var reader = new FileReader();
    reader.onloadend = function () {
        socket.emit('chat image', file.name, reader.result);
        $('#messages').append('<p><span class="text-muted">' + username + ':</span> <img class="msg-img" alt="' + file.name + '" title="' + file.name + '" src="' + reader.result + '"></img></p>');
    }
    if (file) {
        reader.readAsDataURL(file);
    }
});
$('#name-btn').click(function(){
    setName();
});
$('#name-input').keypress(function (e) {
    var key = e.which;
    if(key == 13) {
        setName();
        return false;
    }
});

$('#msg-btn').click(function(){
    sendMsg();
});
$('#msg-input').keydown(function (e) {
    var key = e.which;
    if (key === 13) {
        sendMsg();
        return false;
    } else if (key === 38) {
        if (numPresses > 0) {
            numPresses--;
        }
        $('#msg-input').val(messages[numPresses]);
        return false;
    } else if (key === 40) {
        if (numPresses !== messages.length) {
            numPresses++;
            $('#msg-input').val(messages[numPresses]);
        } else {
            $('#msg-input').val("");
        }
        return false;
    }
});

$("#img-btn").click(function(){
    $("#img-input").click();
});




socket.on('chat nick', function(user){
    $('#messages').append('<p><b>' + user + ' has joined.</b></p>');
    var notification = new Notification(user + " has joined.");
    notification.onshow = function () {
        setTimeout(function () {
            notification.close()
        }, 3000)
    };
});
socket.on('chat setnick', function(user, set){
    username = user;
    if (set) {
        return;
    }
    $('#name-div').addClass('hidden');
    $('#messages').removeClass('hidden');
    $('#msg-input').removeAttr('disabled');
    $('#msg-btn').removeClass('disabled');
    $('#img-btn').removeClass('disabled');
    $('#messages').append('<p><b>' + username + ' has joined.</b></p>');
});
socket.on('chat message', function(msg, user){
    $('#messages').append('<p><span class="text-muted">' + user + ':</span> ' + msg + '</p>');
    var notification = new Notification(user + ": " + msg);
    notification.onshow = function () {
        setTimeout(function () {
            notification.close()
        }, 3000)
    };
});
socket.on('chat command', function(msg, command){
    var comm = command.split(" ")[0];
    $('#messages').append('<p><span class="text-muted"> ' + comm + ':</span> ' + msg + '</p>');
    var notification = new Notification(command + ": " + msg);
    notification.onshow = function () {
        setTimeout(function () {
            notification.close()
        }, 3000)
    };
});
socket.on('chat dm', function(msg, user){
    $('#messages').append('<p><span class="text-muted">DM from ' + user + ':</span> ' + msg + '</p>');
    var notification = new Notification('DM from ' + user + ": " + msg);
    notification.onshow = function () {
        setTimeout(function () {
            notification.close()
        }, 3000)
    };
});
socket.on('chat error', function(errNo){
    if (errNo === 0) {
        alert("That username is already taken");
    } else if (errNo === 1) {
        $('#messages').append('<p><span class="text-muted" style="color: red;"> Error: </span> Your direct message failed.</p>');
    }
});
socket.on('chat image', function(imgName, imgSrc, user){
    $('#messages').append('<p><span class="text-muted">' + user + ':</span> <img class="msg-img" alt="' + imgName + '" title="' + imgName + '" src="' + imgSrc + '"></img></p>');
    var notification = new Notification(user + " has uploaded a photo.");
    notification.onshow = function () {
        setTimeout(function () {
            notification.close()
        }, 3000)
    };
});
socket.on('chat disconnect', function(user){
    $('#messages').append('<p><b>' + user + ' has disconnected.</b></p>');
    var notification = new Notification(user + " has disconnected.");
    notification.onshow = function () {
        setTimeout(function () {
            notification.close()
        }, 3000)
    };
});

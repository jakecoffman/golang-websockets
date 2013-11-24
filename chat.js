$(function () {

	var conn;
	var msg = $("#msg");
	var log = $("#log");

	function appendLog(msg) {
		var d = log[0];
		var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
		msg.appendTo(log);
		if (doScroll) {
			d.scrollTop = d.scrollHeight - d.clientHeight;
		}
	}

	$("#form").submit(function () {
		if (!conn) {
			return false;
		}
		if (!msg.val()) {
			return false;
		}
		conn.send(msg.val());
		msg.val("");
		return false
	});

	if (window["WebSocket"]) {
		conn = new WebSocket("ws://localhost:8080/ws");
		conn.onclose = function(e) {
			appendLog($("<div><b>Connection closed.</b></div>"))
		};

		conn.onmessage = function(e) {
			appendLog($("<div/>").text(e.data))
		};

		conn.onopen = function(e) {
			appendLog($("<div><b>Welcome to the chat.</b></div>"));
		};
	} else {
		appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
	}
});

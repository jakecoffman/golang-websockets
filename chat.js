var app = angular.module("chat", []);

app.controller("MainCtl", function ($scope) {
	$scope.log = [];
	$scope.message = "";

	if (!window["WebSocket"]) {
		$scope.log.push("Your browser does not support WebSockets.");
		return;
	}
	var conn = new WebSocket("ws://localhost:8080/ws");
	conn.onclose = function (e) {
		$scope.$apply(function () {
			$scope.log.push("Connection closed.");
		})
	};

	conn.onmessage = function (e) {
		$scope.$apply(function () {
			$scope.log.push(e.data);
		})
	};

	conn.onopen = function (e) {
		console.log("Connected");
		$scope.$apply(function () {
			$scope.log.push("Welcome to the chat!");
		})
	};

	$scope.send = function () {
		if (!conn) {
			return;
		}

		if (!$scope.message) {
			return;
		}

		conn.send($scope.message);
		$scope.message = "";
	}
});


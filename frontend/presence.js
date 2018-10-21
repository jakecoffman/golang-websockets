var app = angular.module("presence", []);

app.directive("presence", function($location, $anchorScroll){
	return {
		link: function(scope, element, attrs){
			$location.hash('bottom');
			scope.$watch("log", function(){
				$anchorScroll();
			}, true);
		}
	}
});

app.controller("MainCtl", function ($scope) {
	$scope.log = [];
	$scope.users = [];

	$scope.room = prompt("Enter room:");
	$scope.user = prompt("Enter user:");

	if (!window["WebSocket"]) {
		$scope.log.push("Your browser does not support WebSockets.");
		return;
	}

	// Connect
	var conn = new WebSocket("ws://localhost:8080/ws?room=" + $scope.room + "&user=" + $scope.user);

	conn.onclose = function (e) {
		$scope.$apply(function () {
			$scope.log.push("Connection closed.");
		})
	};

	conn.onmessage = function (e) {
		$scope.$apply(function () {
			var data = JSON.parse(e.data);
			console.log(data);
			$scope.users = data.users
			$scope.log.push(data.action);
		})
	};

	conn.onopen = function (e) {
		console.log("Connected");
		$scope.$apply(function () {
			$scope.log.push("Welcome to the presence example!");
		})
	};
});

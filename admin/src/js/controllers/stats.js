angular.module('RDash')
	.controller('StatsCtrl', ['$scope', '$http','$modal', statsCtrl]);

function statsCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Stats list");
	$scope.stats = [];

	$http.post(APIURL + '/getStatsList', {sessionid: $scope.Self.SessionId})
		.success(function (data) {
			$scope.stats = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch stats infos!", "danger");
		});

	$scope.addStat = function () {
		$scope.goto("/achievementAdd");
	}
}
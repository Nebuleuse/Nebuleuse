angular.module('RDash')
	.controller('UpdatesCtrl', ['$scope', '$http','$modal', UpdatesCtrl]);

function UpdatesCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Achievements list");
	if(!$scope.checkAccess())
		return;
	$scope.list = [];

	$http.post(APIURL + '/getUpdateListWithGit', {sessionid: $scope.Self.SessionId})
		.success(function (data) {
			$scope.list = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch updates infos!", "danger");
		});

	$scope.addAchievement = function () {
		$scope.goto("/updateAdd");
	}
}
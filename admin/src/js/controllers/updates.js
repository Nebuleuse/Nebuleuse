angular.module('RDash')
	.controller('UpdatesCtrl', ['$scope', '$http','$uibModal', UpdatesCtrl]);

function UpdatesCtrl($scope, $http, $uibModal) {
	$scope.setPageTitle("Update list");
	if(!$scope.checkAccess())
		return;
	$scope.list = [];

	$http.post(APIURL + '/getUpdateListWithGit', {sessionid: $scope.Self.SessionId, diffs: true})
		.success(function (data) {
			console.dir(data)
			$scope.list = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch updates infos!", "danger");
		});

	$scope.addAchievement = function () {
		$scope.goto("/updateAdd");
	}
}
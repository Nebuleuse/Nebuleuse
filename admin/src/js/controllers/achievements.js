angular.module('RDash')
	.controller('AchievementsCtrl', ['$scope', '$http','$modal', AchievementsCtrl]);

function AchievementsCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Achievements list");
	if(!$scope.checkAccess())
		return;
	$scope.achievements = [];

	$http.post(APIURL + '/getAchievements', {sessionid: $scope.Self.SessionId})
		.success(function (data) {
			$scope.achievements = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch achievements infos!", "danger");
		});

	$scope.addAchievement = function () {
		$scope.goto("/achievementAdd");
	}
}
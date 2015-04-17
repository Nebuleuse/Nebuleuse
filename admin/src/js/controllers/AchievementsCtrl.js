angular.module('RDash')
    .controller('AchievementsCtrl', ['$scope', '$http', AchievementsCtrl]);

function AchievementsCtrl($scope, $http) {
	$scope.setPageTitle("Achievements");
	$scope.achievements = [];
	$http.post(APIURL + '/getAchievements', {sessionid: $scope.Self.SessionId})
		.success(function (data) {
			$scope.achievements = data.Achievements;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch achievements infos!", "danger");
		});
}
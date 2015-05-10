angular.module('RDash')
	.controller('AchievementCtrl', ['$scope', '$http', '$modal', '$location', '$stateParams', AchievementCtrl]);

function AchievementCtrl($scope, $http, $modal, $location, $stateParams) {
	$scope.getAchievement = function () {
		$http.post(APIURL + '/getAchievement', {sessionid: $scope.Self.SessionId, achievementid: $stateParams.achievementId})
		.success(function (data) {
			$scope.achievement = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch achievement infos!", "danger");
		});
	}

	$scope.setPageTitle("Achievements info");
	$scope.achievement = {};
	$scope.editing = false;

	if ($location.path().startsWith("/achievementEdit")) {
		$scope.getAchievement();
	} else if ($location.path().startsWith("/achievementAdd")) {
		$scope.achievement.Id == null;
	}

	$scope.startEdit = function () {
		$scope.editing = true;
	}
	$scope.cancelEdit = function () {
		$scope.editing = false;
		$scope.getAchievement();
	}
	$scope.saveEdit = function () {
		var ach = $scope.achievement;
		if(ach.Id === null)
			return $scope.saveAchievement();

		var toSend = {sessionid: $scope.Self.SessionId, achievementid: ach.Id, data: angular.toJson(ach)};
		$http.post(APIURL + '/setAchievement', toSend)
		.success(function (data) {
			$scope.editing = false;
		}).error(function (data, status) {
			$scope.addAlert("Could not save achievement infos!", "danger");
		});
	}
	$scope.saveAchievement = function() {
		var ach = $scope.achievement;
		var toSend = {sessionid: $scope.Self.SessionId, data: angular.toJson(ach)};
		$http.post(APIURL + '/addAchievement', toSend)
		.success(function (data) {
			$scope.editing = false;
			ach.Id = data.Id
		}).error(function (data, status) {
			$scope.addAlert("Could not add achievement!", "danger");
		});
	};
	$scope.deleteAchievement = function () {
		var modalInstance = $modal.open({
			templateUrl: 'templates/confirmDelete.html',
			controller: 'ModalCtrl'
		});

		modalInstance.result.then(function () {
			var toSend = {sessionid: $scope.Self.SessionId, achievementid: $scope.achievement.Id };
			$http.post(APIURL + '/deleteAchievement', toSend)
			.success(function (data) {
				$scope.goto("/achievements")
			}).error(function (data, status) {
				$scope.addAlert("Could not delete achievement!", "danger");
			});
		}, function () {
			return;
		});
	}
}
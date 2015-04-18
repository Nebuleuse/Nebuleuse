angular.module('RDash')
	.controller('AchievementsCtrl', ['$scope', '$http','$modal', AchievementsCtrl]);

function AchievementsCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Achievements");
	$scope.achievements = [];
	$scope.editing = -1;
	$scope.originAchievement = {};

	$http.post(APIURL + '/getAchievements', {sessionid: $scope.Self.SessionId})
		.success(function (data) {
			$scope.achievements = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch achievements infos!", "danger");
		});

	$scope.selectEdit  = function(index) {
		if($scope.editing != -1)
			$scope.cancelEdit();

		$scope.originAchievement = angular.toJson($scope.achievements[index], false);
		$scope.editing = index;
	};
	$scope.saveEdit = function (index) {
		var current = $scope.achievements[index];
		if(current.Id == null)
			return $scope.saveAchievement(current);

		var toSend = {sessionid: $scope.Self.SessionId, achievementid: current.Id, data: angular.toJson(current)};
		$http.post(APIURL + '/setAchievement', toSend)
		.success(function (data) {
			$scope.editing = -1;
		}).error(function (data, status) {
			$scope.cancelEdit();
			$scope.addAlert("Could not save achievements infos!", "danger");
		});
	}
	$scope.saveAchievement = function(ach) {
		var toSend = {sessionid: $scope.Self.SessionId, data: angular.toJson(ach)};
		$http.post(APIURL + '/addAchievement', toSend)
		.success(function (data) {
			ach.Id = data.Id
			$scope.editing = -1;
		}).error(function (data, status) {
			$scope.addAlert("Could not add achievements!", "danger");
		});
	};
	$scope.cancelEdit = function () {
		if($scope.achievements[$scope.editing].Id == null){
			$scope.achievements.splice($scope.editing, 1);
		} else {
			$scope.achievements[$scope.editing] = angular.fromJson($scope.originAchievement);
		}
		$scope.editing = -1;
	}
	$scope.addAchievement = function () {
		$scope.editing = $scope.achievements.length;
		$scope.achievements.push({});
	}
	$scope.deleteAchievement = function (index) {
		var modalInstance = $modal.open({
			templateUrl: 'templates/confirmDelete.html',
			controller: 'ModalCtrl'
		});

		modalInstance.result.then(function () {
			var toSend = {sessionid: $scope.Self.SessionId, achievementid: $scope.achievements[index].Id };
			$http.post(APIURL + '/deleteAchievement', toSend)
			.success(function (data) {
				$scope.achievements.splice(index, 1);
			}).error(function (data, status) {
				$scope.addAlert("Could not delete achievement!", "danger");
			});
		}, function () {
			return;
		});
	}
}
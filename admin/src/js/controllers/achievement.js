angular.module('RDash')
	.controller('AchievementCtrl', ['$scope', '$http','$modal', AchievementCtrl]);

function AchievementCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Achievements info");
	$scope.achievement = {};

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
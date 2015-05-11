angular.module('RDash')
	.controller('StatsCtrl', ['$scope', '$http','$modal', statsCtrl]);

function statsCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Stats list");
	if(!$scope.checkAccess())
		return;
	$scope.stats = [];

	$http.post(APIURL + '/getStatsList', {sessionid: $scope.Self.SessionId})
		.success(function (data) {
			$scope.stats = data;
			$scope.usersStats = {Fields: [], ExtraFields:[]};
			for (var i = $scope.stats.length - 1; i >= 0; i--) {
				if($scope.stats[i].Name == 'users')
					$scope.usersStats.Fields = $scope.usersStats.Fields.concat($scope.stats[i].Fields);
				else if ($scope.stats[i].AutoCount)
					$scope.usersStats.ExtraFields.push($scope.stats[i].Name);
		};
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch stats infos!", "danger");
		});

	$scope.addStat = function () {
		$scope.goto("/achievementAdd");
	}
}
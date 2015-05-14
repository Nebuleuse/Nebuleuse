angular.module('RDash')
	.controller('StatsCtrl', ['$scope', '$http','$modal', statsCtrl]);

function statsCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Stat tables list");
	if(!$scope.checkAccess())
		return;
	$scope.stats = [];

	$http.post(APIURL + '/getStatTables', {sessionid: $scope.Self.SessionId})
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

	$scope.getFields = function (Fields) {
		var ret = "";
		for (var i = Fields.length - 1; i >= 0; i--) {
			ret += Fields[i].Name + " ";
		};
		return ret;
	}
}
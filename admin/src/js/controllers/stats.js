angular.module('RDash')
	.controller('StatsCtrl', ['$scope', '$http','$modal', statsCtrl]);

function statsCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Stat tables list");
	if(!$scope.checkAccess())
		return;

	$scope.editing = false;
	$scope.stats = [];

	$scope.getStatsTables = function () {
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
	}
	$scope.getStatsTables();

	$scope.addUsersField = function () {
		$scope.usersStats.Fields[$scope.usersStats.Fields.length] = {};
	}
	$scope.removeUsersField = function (field) {
		for (var i = $scope.usersStats.Fields.length - 1; i >= 0; i--) {
			if ($scope.usersStats.Fields[i] === field){
				$scope.usersStats.Fields.splice(i, 1);
			}
		};
	}
	$scope.saveUsersFields = function () {
		var fields = "";
		for (var i = $scope.usersStats.Fields.length - 1; i >= 0; i--) {
			if(i == 0)
				fields += $scope.usersStats.Fields[i].Name
			else
				fields += $scope.usersStats.Fields[i].Name + ","
		};
		$http.post(APIURL + '/setUsersStatFields', {sessionid: $scope.Self.SessionId, fields: fields})
			.success(function (data) {
				$scope.editing = false;
			}).error(function (data, status) {
				$scope.parseError(data, status);
				$scope.addAlert("Could not set users stats fields!", "danger");
			})
		console.log(fields);
	}

	$scope.getFields = function (Fields) {
		var ret = "";
		for (var i = Fields.length - 1; i >= 0; i--) {
			ret += Fields[i].Name + " ";
		};
		return ret;
	}

	$scope.startEdit = function() {
		$scope.editing = true;
	};
	$scope.cancelEdit = function() {
		$scope.editing = false;
		$scope.getStatsTables();
	};
}
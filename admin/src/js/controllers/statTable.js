angular.module('RDash')
	.controller('StatTableEditCtrl', ['$scope', '$http', '$uibModal', '$location', '$stateParams', StatTableEditCtrl]);

function StatTableEditCtrl($scope, $http, $modal, $location, $stateParams) {
	$scope.getStatTable = function () {
		$http.post(APIURL + '/getStatTable', {sessionid: $scope.Self.SessionId, name: $scope.table.Name})
		.success(function (data) {
			$scope.table = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch table infos!", "danger");
		});
	}

	$scope.setPageTitle("Stat table infos");
	if(!$scope.checkAccess())
		return;
	$scope.table = {};
	$scope.editing = false;

	if ($location.path().startsWith("/statTableEdit")) {
		$scope.table.Name = $stateParams.statName;
		$scope.getStatTable();
	} else if ($location.path().startsWith("/statTableAdd")) {
		$scope.toAdd = true;
		$scope.table.Name = "";
		$scope.table.Fields = [];
		$scope.table.AutoCount = false;
		$scope.editing = true;
	}

	$scope.startEdit = function () {
		$scope.editing = true;
	}
	$scope.cancelEdit = function () {
		$scope.editing = false;
		if($scope.toAdd)
			return $scope.goto('/stats');
		$scope.getStatTable();
	}
	$scope.addField = function () {
		$scope.table.Fields[$scope.table.Fields.length] = {};
	}
	$scope.removeField = function (index) {
		$scope.table.Fields.splice(index, 1);
	}
	$scope.saveEdit = function () {
		var table = $scope.table;
		if($scope.toAdd)
			return $scope.saveTable();

		var toSend = {sessionid: $scope.Self.SessionId, data: angular.toJson(table)};
		$http.post(APIURL + '/setStatTable', toSend)
		.success(function (data) {
			$scope.editing = false;
		}).error(function (data, status) {
			$scope.addAlert("Could not save table infos!", "danger");
		});
	}
	$scope.saveTable = function() {
		var toSend = {sessionid: $scope.Self.SessionId, data: angular.toJson($scope.table)};
		console.log(toSend);
		$http.post(APIURL + '/addStatTable', toSend)
		.success(function (data) {
			$scope.editing = false;
			$scope.toAdd = false;
		}).error(function (data, status) {
			$scope.addAlert("Could not add Table!", "danger");
		});
	};
	$scope.deleteTable = function () {
		var modalInstance = $modal.open({
			templateUrl: 'templates/confirmDelete.html',
			controller: 'ModalCtrl'
		});

		modalInstance.result.then(function () {
			var toSend = {sessionid: $scope.Self.SessionId, name: $scope.table.Name };
			$http.post(APIURL + '/deleteStatTable', toSend)
			.success(function (data) {
				$scope.goto("/stats")
			}).error(function (data, status) {
				$scope.addAlert("Could not delete table!", "danger");
			});
		}, function () {
			return;
		});
	}
}
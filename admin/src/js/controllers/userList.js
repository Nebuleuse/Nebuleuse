angular.module('RDash')
    .controller('UserListCtrl', ['$scope', '$http', UserListCtrl]);

function UserListCtrl($scope, $http) {
	$scope.users= [];
	$scope.page = 1;

	$scope.pageChanged = function () {
		$http.post(APIURL + '/getUsersInfos', {sessionid: $scope.Self.SessionId, page: $scope.page, infomask: 1})
		.success(function (data) {
			$scope.users = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch users infos!", "danger");
		});
	}
	$scope.pageChanged();
}

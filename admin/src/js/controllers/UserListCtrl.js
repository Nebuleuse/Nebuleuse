angular.module('RDash')
    .controller('UserListCtrl', ['$scope', '$cookieStore','$http', '$location', UserListCtrl]);

function UserListCtrl($scope, $cookieStore, $http, $location) {
	$scope.users= [];
	$scope.page = 1;

	$http.post(APIURL + '/getUsersInfos', {sessionid: $scope.Self.SessionId, page: $scope.page, infomask: 1})
	.success(function (data) {
		$scope.users = data;
	}).error(function (data) {
		$scope.addAlert("Could not fetch users infos!", "error");
	});
}

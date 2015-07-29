angular.module('RDash')
    .controller('UserListCtrl', ['$scope', '$http', "$attrs", UserListCtrl]);

function UserListCtrl($scope, $http, $attrs) {
	if($attrs.changetitle != "false")
		$scope.setPageTitle("Users list");

	if(!$scope.checkAccess())
		return;
	
	$scope.users= [];
	$scope.page = 1;

	$scope.pageChanged = function () {
		console.log({sessionid: $scope.Self.SessionId, infomask:UserMaskBase});
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

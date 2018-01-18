angular.module('RDash')
    .controller('UserCtrl', ['$scope', '$http', '$location', '$stateParams', UserCtrl]);

function UserCtrl($scope, $http, $location, $stateParams) {
	$scope.setPageTitle("User infos");
	if(!$scope.checkAccess())
		return;
	
	$scope.getUser = function () {
		$http.post(APIURL + '/getUserInfos', {sessionid: $scope.Self.SessionId, userid: $scope.user.Id, infomask: UserMaskAll})
		.success(function (data) {
			$scope.user = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch user infos!", "danger");
		});
	}

	$scope.user = {};
	$scope.editing = false;

	if ($location.path().startsWith("/user")) {
		$scope.user.Id = $stateParams.userId;
		$scope.getUser();
	} else if ($location.path().startsWith("/userAdd")) {
		$scope.toAdd = true;
		$scope.editing = true;
	}

		
}

angular.module('RDash')
    .controller('OnlineUserListCtrl', ['$scope', '$http', "$attrs", OnlineUserListCtrl]);

function OnlineUserListCtrl($scope, $http, $attrs) {
	if($attrs.changetitle != "false")
		$scope.setPageTitle("Users list");

	if(!$scope.checkAccess())
		return;
	
	$scope.users= [];
	$scope.page = 1;

	$http.post(APIURL + '/getOnlineUsersList', {sessionid: $scope.Self.SessionId})
	.success(function (data) {
		$scope.users = data;
	}).error(function (data, status) {
		$scope.parseError(data, status);
		$scope.addAlert("Could not fetch online users infos!", "danger");
	});
}

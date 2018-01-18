angular.module('RDash')
    .controller('DashboardCtrl', ['$scope', '$cookieStore','$http', '$location', DashboardCtrl]);

function DashboardCtrl($scope, $cookieStore, $http, $location) {
	$scope.setPageTitle("Dashboard");
	if(!$scope.checkAccess())
		return;

	$scope.infos = {};

	$http.post(APIURL + '/getDashboardInfos', {sessionid: $scope.Self.SessionId})
	.success(function (data) {
		$scope.infos = data;
	}).error(function (data) {
		$scope.parseError(data, status);
		$scope.addAlert("Could not fetch dashboard infos!", "danger");
	});
}
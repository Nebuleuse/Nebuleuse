angular.module('RDash')
    .controller('DashboardCtrl', ['$scope', '$cookieStore','$http', '$location', DashboardCtrl]);

function DashboardCtrl($scope, $cookieStore, $http, $location) {
	$scope.setPageTitle("Dashboard");
	$scope.checkAccess();
}
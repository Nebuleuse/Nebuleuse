angular.module('RDash')
	.controller('InstallCtrl', ['$scope', '$http','$uibModal', InstallCtrl]);

function InstallCtrl($scope, $http, $modal) {
	$scope.setPageTitle("Installation");
}
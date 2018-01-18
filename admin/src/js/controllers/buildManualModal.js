angular.module('RDash')
	.controller('ManualBulidCreateModal', ['$scope', '$http', '$uibModalInstance', ManualBulidCreateModal]);

function ManualBulidCreateModal($scope, $http, $uibModalInstance) {
	$scope.log = "";
    $scope.close = function(){
        $uibModalInstance.close();
    }
    $scope.createManualBuild = function(log){
        //todo
       $http.post(APIURL + '/addBuild', {sessionid: $scope.Self.SessionId, "log": log})
		.success(function () {
			$scope.refreshList();
			$uibModalInstance.close();
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not create update!", "danger");
		});
    }
}
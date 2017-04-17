angular.module('RDash')
	.controller('BranchCreateModal', ['$scope', '$http', '$uibModalInstance', 'build', BranchCreateModal]);

function BranchCreateModal($scope, $http, $uibModalInstance, build) {
    $scope.build = build;
    $scope.empty = build == 0;
    $scope.close = function(){
        $uibModalInstance.close();
    }
    $scope.addBranch = function(name, accessrank, log, semver){
        if($scope.empty){
            $http.post(APIURL + '/addEmptyBranch', {sessionid: $scope.Self.SessionId, "name": name, "accessrank": accessrank})
            .success(function () {
                $scope.refreshList();
                $uibModalInstance.close();
            }).error(function (data, status) {
                $scope.parseError(data, status);
                $scope.addAlert("Could not create branch!", "danger");
            });
        } else {
            $http.post(APIURL + '/addBranchFromBuild', {sessionid: $scope.Self.SessionId, "name": name, "accessrank": accessrank, "semver": semver, "log": log, "build": $scope.build})
            .success(function () {
                $scope.refreshList();
                $uibModalInstance.close();
            }).error(function (data, status) {
                $scope.parseError(data, status);
                $scope.addAlert("Could not create branch!", "danger");
            });
        }
       
    }
}
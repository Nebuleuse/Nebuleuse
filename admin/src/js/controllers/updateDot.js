angular.module('RDash')
	.controller('updateDotCtrl', ['$scope', updateDotCtrl]);

function updateDotCtrl($scope) {
    var updates = $scope.branch.Updates;
    $scope.update = {build: $scope.build, branch: $scope.branch, create: true};
    $scope.exist = false;
    $scope.visible = true;
    $scope.active = false;
	for (var i=0; i < updates.length; i++){
		if(updates[i].BuildId == $scope.build.Id){
			$scope.update = updates[i];
            $scope.exist = true;
            if($scope.branch.ActiveBuild == $scope.build.Id){
                $scope.active = true;
            }
		}
	}
    if(!$scope.exist){
        if($scope.branch.Updates[0].BuildId >= $scope.build.Id){
            $scope.visible = false;
        }
    }
}